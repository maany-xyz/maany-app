package keeper

import (
    "context"
    "fmt"
    "time"

    "cosmossdk.io/log"
    storetypes "cosmossdk.io/store/types"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
    paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
    "github.com/maany-xyz/maany-app/x/autolp/types"
    icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
    icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
    icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
    channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
    clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
    wraptypes "github.com/maany-xyz/maany-app/x/transfer/types"
    "strings"
)

type Keeper struct {
    cdc      codec.BinaryCodec
    storeKey storetypes.StoreKey
    ps       paramtypes.Subspace

    transferKeeper    types.TransferKeeper
    interchainQuery   types.InterchainTxsQueryKeeper

    // ICA controller
    icaControllerKeeper    icacontrollerkeeper.Keeper
    icaControllerMsgServer icacontrollertypes.MsgServer

    // authority address (e.g., gov/admin module address)
    authority string
}

// lastAttemptHeightKey builds the store key to track last handshake attempt
// height for a (connectionID, icaID) pair.
func lastAttemptHeightKey(connectionID, icaID string) []byte {
    return []byte(fmt.Sprintf("last_ica_attempt/%s/%s", connectionID, icaID))
}

// pendingChanKey tracks a pending ICA channel handshake per owner
// (owner is the raw owner string used to derive controller port ID).
func pendingChanKey(owner string) []byte {
    return []byte(fmt.Sprintf("pending_ica_channel/%s", owner))
}

func NewKeeper(
    cdc codec.BinaryCodec,
    key storetypes.StoreKey,
    ps paramtypes.Subspace,
    transferKeeper types.TransferKeeper,
    interchainQuery types.InterchainTxsQueryKeeper,
    icaControllerKeeper icacontrollerkeeper.Keeper,
    icaControllerMsgServer icacontrollertypes.MsgServer,
    authority string,
) Keeper {
    if !ps.HasKeyTable() {
        ps = ps.WithKeyTable(types.ParamKeyTable())
    }
    return Keeper{
        cdc:            cdc,
        storeKey:       key,
        ps:             ps,
        transferKeeper: transferKeeper,
        interchainQuery: interchainQuery,
        icaControllerKeeper:    icaControllerKeeper,
        icaControllerMsgServer: icaControllerMsgServer,
        authority:              authority,
    }
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
    return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
    var p types.Params
    k.ps.GetParamSet(ctx, &p)
    return p
}

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) { k.ps.SetParamSet(ctx, &p) }

// Query server passthrough
var _ types.QueryServer = Keeper{}

func (k Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
    return &types.QueryParamsResponse{Params: types.Params{}}, nil
}

func (k Keeper) InterchainAccountAddress(c context.Context, req *types.QueryICAAddressRequest) (*types.QueryICAAddressResponse, error) {
    if req == nil {
        return nil, fmt.Errorf("invalid request")
    }
    ctx := sdk.UnwrapSDKContext(c)
    owner := k.ownerString(req.InterchainAccountId)
    portID, err := icatypes.NewControllerPortID(owner)
    if err != nil {
        return nil, err
    }
    addr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, req.ConnectionId, portID)
    if !found {
        return nil, fmt.Errorf("no interchain account found for port %s", portID)
    }
    return &types.QueryICAAddressResponse{InterchainAccountAddress: addr}, nil
}

// ownerString builds the ICA owner identifier used to derive the controller port ID.
// Must not contain '/' since it's used in an IBC identifier. Use a dot as separator.
// ownerString returns the bech32 autolp module account address.
// ICA owner must be strictly alphanumeric; a bech32 account address satisfies this.
// We use a single owner for autolp-managed ICA on a given connection.
func (k Keeper) ownerString(_ string) string {
    return authtypes.NewModuleAddress(types.ModuleName).String()
}

// IBC event handlers (minimal): emit events for acks/timeouts/open-ack
func (k Keeper) HandleChanOpenAck(ctx sdk.Context, portID, channelID, counterpartyChannelID, counterpartyVersion string) error {
    k.Logger(ctx).Info("ICA channel open ack", "port_id", portID, "channel_id", channelID, "cp_channel_id", counterpartyChannelID, "cp_version", counterpartyVersion)
    // Clear pending flag for this owner (owner is suffix after controller prefix).
    owner := portID
    if strings.HasPrefix(owner, icatypes.ControllerPortPrefix) {
        owner = strings.TrimPrefix(owner, icatypes.ControllerPortPrefix)
    }
    store := ctx.KVStore(k.storeKey)
    store.Delete(pendingChanKey(owner))
    ctx.EventManager().EmitEvent(sdk.NewEvent(
        "autolp_ica_chan_open_ack",
        sdk.NewAttribute("port_id", portID),
        sdk.NewAttribute("channel_id", channelID),
        sdk.NewAttribute("cp_channel_id", counterpartyChannelID),
        sdk.NewAttribute("cp_version", counterpartyVersion),
    ))
    return nil
}

func (k Keeper) HandleAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, relayer sdk.AccAddress) error {
    k.Logger(ctx).Info("ICA packet ack", "src_port", packet.SourcePort, "src_channel", packet.SourceChannel, "sequence", packet.Sequence, "relayer", relayer.String())
    ctx.EventManager().EmitEvent(sdk.NewEvent(
        "autolp_ica_ack",
        sdk.NewAttribute("src_port", packet.SourcePort),
        sdk.NewAttribute("src_channel", packet.SourceChannel),
        sdk.NewAttribute("sequence", fmt.Sprintf("%d", packet.Sequence)),
        sdk.NewAttribute("relayer", relayer.String()),
    ))
    return nil
}

func (k Keeper) HandleTimeout(ctx sdk.Context, packet channeltypes.Packet, relayer sdk.AccAddress) error {
    k.Logger(ctx).Warn("ICA packet timeout", "src_port", packet.SourcePort, "src_channel", packet.SourceChannel, "sequence", packet.Sequence, "relayer", relayer.String())
    ctx.EventManager().EmitEvent(sdk.NewEvent(
        "autolp_ica_timeout",
        sdk.NewAttribute("src_port", packet.SourcePort),
        sdk.NewAttribute("src_channel", packet.SourceChannel),
        sdk.NewAttribute("sequence", fmt.Sprintf("%d", packet.Sequence)),
        sdk.NewAttribute("relayer", relayer.String()),
    ))
    return nil
}

// SendAppToICA sends the provided amount from the autolp module account to the
// ICA address on the host chain over the provided ICS-20 transfer channel.
// This is an internal helper; no signatures required.
func (k Keeper) SendAppToICA(ctx sdk.Context, connectionID, icaID, transferChannel string, amount sdk.Coin, timeoutSec uint64) error {
    // Resolve ICA address on host chain
    owner := k.ownerString(icaID)
    portID, err := icatypes.NewControllerPortID(owner)
    if err != nil {
        return err
    }
    icaAddr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, connectionID, portID)
    if !found || icaAddr == "" {
        k.Logger(ctx).Warn("ICA address not found; cannot fund", "port_id", portID, "connection_id", connectionID)
        return fmt.Errorf("ica address not found for %s on %s", portID, connectionID)
    }

    // Build transfer from module account to ICA address
    sender := authtypes.NewModuleAddress(types.ModuleName).String()
    msg := &wraptypes.MsgTransfer{
        SourcePort:       "transfer",
        SourceChannel:    transferChannel,
        Token:            amount,
        Sender:           sender,
        Receiver:         icaAddr,
        TimeoutHeight:    clienttypes.Height{},
        TimeoutTimestamp: uint64(time.Now().Add(time.Duration(timeoutSec) * time.Second).UnixNano()),
        Memo:             "autolp:fund-ica",
    }

    k.Logger(ctx).Info("Funding ICA via ICS-20", "connection_id", connectionID, "ica_id", icaID, "channel", transferChannel, "amount", amount.String(), "receiver", icaAddr)
    // Call wrapped transfer keeper
    _, err = k.transferKeeper.Transfer(sdk.WrapSDKContext(ctx), msg)
    return err
}

// EnsureICA ensures there is an active ICA controller channel for the given
// (connectionID, icaID). If not present, it initiates the ICA registration
// handshake by calling the ICA controller MsgServer. This is safe to call
// from Begin/EndBlocker or cron; it is idempotent as long as the channel is
// already active or once the handshake completes.
func (k Keeper) EnsureICA(ctx sdk.Context, connectionID, icaID string) error {
    owner := k.ownerString(icaID)
    portID, err := icatypes.NewControllerPortID(owner)
    if err != nil {
        return err
    }
    if _, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionID, portID); found {
        k.Logger(ctx).Info("ICA already active", "connection_id", connectionID, "ica_id", icaID, "port_id", portID)
        return nil
    }
    // Initiate registration (ORDERED is standard for ICA)
    k.Logger(ctx).Info("Initiating ICA registration", "connection_id", connectionID, "ica_id", icaID, "port_id", portID)
    resp, err := k.icaControllerMsgServer.RegisterInterchainAccount(ctx, &icacontrollertypes.MsgRegisterInterchainAccount{
        Owner:        owner,
        ConnectionId: connectionID,
        Version:      "",
        Ordering:     channeltypes.ORDERED,
    })
    if err != nil {
        return err
    }
    // Mark pending handshake for this owner to avoid spamming new channels.
    store := ctx.KVStore(k.storeKey)
    store.Set(pendingChanKey(owner), []byte(resp.ChannelId))
    return nil
}

// EnsureICAIfDue reads Params and attempts to initiate ICA registration every
// EnsureIcaInterval blocks until a channel becomes active.
func (k Keeper) EnsureICAIfDue(ctx sdk.Context) {
    params := k.GetParams(ctx)
    if !params.EnsureIcaEnabled {
        k.Logger(ctx).Debug("EnsureICA disabled by params")
        return
    }
    if params.ConnectionId == "" || params.IcaId == "" {
        k.Logger(ctx).Warn("EnsureICA misconfigured: missing connection_id or ica_id", "connection_id", params.ConnectionId, "ica_id", params.IcaId)
        return
    }
    owner := k.ownerString(params.IcaId)
    portID, err := icatypes.NewControllerPortID(owner)
    if err != nil {
        k.Logger(ctx).Error("Failed to derive controller portID", "owner", owner, "err", err)
        return
    }
    // If we have a pending handshake recorded, wait for ack instead of opening new channels.
    if ctx.KVStore(k.storeKey).Get(pendingChanKey(owner)) != nil {
        k.Logger(ctx).Debug("EnsureICA: pending handshake exists, waiting for ack", "owner", owner)
        return
    }
    if _, found := k.icaControllerKeeper.GetActiveChannelID(ctx, params.ConnectionId, portID); found {
        k.Logger(ctx).Debug("ICA channel already active", "connection_id", params.ConnectionId, "ica_id", params.IcaId, "port_id", portID)
        return
    }
    // throttle attempts
    store := ctx.KVStore(k.storeKey)
    key := lastAttemptHeightKey(params.ConnectionId, params.IcaId)
    var last int64 = 0
    if bz := store.Get(key); bz != nil {
        last = int64(sdk.BigEndianToUint64(bz))
    }
    if ctx.BlockHeight() < last+int64(params.EnsureIcaInterval) {
        k.Logger(ctx).Debug("EnsureICA backoff in effect", "last_attempt", last, "interval", params.EnsureIcaInterval, "height", ctx.BlockHeight())
        return
    }
    // attempt
    if err := k.EnsureICA(ctx, params.ConnectionId, params.IcaId); err != nil {
        k.Logger(ctx).Error("EnsureICA attempt failed", "connection_id", params.ConnectionId, "ica_id", params.IcaId, "err", err)
    } else {
        k.Logger(ctx).Info("EnsureICA attempt submitted", "connection_id", params.ConnectionId, "ica_id", params.IcaId)
    }
    // record last attempt height regardless of error to back off
    store.Set(key, sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight())))
}
