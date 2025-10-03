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
    icqkeeper "github.com/maany-xyz/maany-app/x/interchainqueries/keeper"
    iqtypes "github.com/maany-xyz/maany-app/x/interchainqueries/types"
    transfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
    icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
    icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
    icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
    channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
    clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
    wraptypes "github.com/maany-xyz/maany-app/x/transfer/types"
    "strings"
    addrcodec "github.com/cosmos/cosmos-sdk/codec/address"
    addrtypes "github.com/cosmos/cosmos-sdk/types/address"
    sdkmath "cosmossdk.io/math"
)

type Keeper struct {
    cdc      codec.BinaryCodec
    storeKey storetypes.StoreKey
    ps       paramtypes.Subspace

    transferKeeper    types.TransferKeeper
    interchainQuery   types.InterchainTxsQueryKeeper
    icqKeeper         *icqkeeper.Keeper

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

// lpStateKey stores the autolp lifecycle state per (connectionID, icaID).
func lpStateKey(connectionID, icaID string) []byte {
    return []byte(fmt.Sprintf("lp_state/%s/%s", connectionID, icaID))
}

// balanceQueryKey stores the ICQ query id for a denom
func balanceQueryKey(connectionID, icaID, hostDenom string) []byte {
    return []byte(fmt.Sprintf("balance_query/%s/%s/%s", connectionID, icaID, hostDenom))
}

func (k Keeper) setBalanceQueryID(ctx sdk.Context, connectionID, icaID, hostDenom string, id uint64) {
    ctx.KVStore(k.storeKey).Set(balanceQueryKey(connectionID, icaID, hostDenom), sdk.Uint64ToBigEndian(id))
}

func (k Keeper) getBalanceQueryID(ctx sdk.Context, connectionID, icaID, hostDenom string) (uint64, bool) {
    bz := ctx.KVStore(k.storeKey).Get(balanceQueryKey(connectionID, icaID, hostDenom))
    if bz == nil {
        return 0, false
    }
    return sdk.BigEndianToUint64(bz), true
}

func NewKeeper(
    cdc codec.BinaryCodec,
    key storetypes.StoreKey,
    ps paramtypes.Subspace,
    transferKeeper types.TransferKeeper,
    interchainQuery types.InterchainTxsQueryKeeper,
    icqKeeper *icqkeeper.Keeper,
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
        icqKeeper:       icqKeeper,
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
    ctx := sdk.UnwrapSDKContext(c)
    return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
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

// LPState implements QueryServer and returns the current lifecycle state.
func (k Keeper) LPState(c context.Context, req *types.QueryLPStateRequest) (*types.QueryLPStateResponse, error) {
    if req == nil {
        return nil, fmt.Errorf("invalid request")
    }
    ctx := sdk.UnwrapSDKContext(c)
    state := k.GetLPState(ctx, req.ConnectionId, req.InterchainAccountId)
    return &types.QueryLPStateResponse{State: string(state)}, nil
}

// SetLPState sets lifecycle state for (connectionID, icaID).
func (k Keeper) SetLPState(ctx sdk.Context, connectionID, icaID string, state types.LPState) {
    prev := k.GetLPState(ctx, connectionID, icaID)
    store := ctx.KVStore(k.storeKey)
    store.Set(lpStateKey(connectionID, icaID), []byte(state))
    if prev != state {
        k.Logger(ctx).Info("LP state transition", "connection_id", connectionID, "ica_id", icaID, "from", string(prev), "to", string(state))
    }
}

// GetLPState returns state, defaulting to UNINITIALIZED if unset.
func (k Keeper) GetLPState(ctx sdk.Context, connectionID, icaID string) types.LPState {
    bz := ctx.KVStore(k.storeKey).Get(lpStateKey(connectionID, icaID))
    if bz == nil {
        return types.LPStateUninitialized
    }
    return types.LPState(string(bz))
}

// Helper transition methods for future automation logic
func (k Keeper) MarkLPReady(ctx sdk.Context, connectionID, icaID string) {
    k.SetLPState(ctx, connectionID, icaID, types.LPStateReady)
}
func (k Keeper) MarkLPCreatingPool(ctx sdk.Context, connectionID, icaID string) {
    k.SetLPState(ctx, connectionID, icaID, types.LPStateCreatingPool)
}
func (k Keeper) MarkLPActive(ctx sdk.Context, connectionID, icaID string) {
    k.SetLPState(ctx, connectionID, icaID, types.LPStateActive)
}
func (k Keeper) PauseLP(ctx sdk.Context, connectionID, icaID string) {
    k.SetLPState(ctx, connectionID, icaID, types.LPStatePaused)
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
    // On successful open-ack, mark state as FUNDING (ICA is active and ready to be funded)
    // Note: owner here is the controller-side owner string derived from port; use params for connection/ica id.
    params := k.GetParams(ctx)
    if params.ConnectionId != "" && params.IcaId != "" {
        k.SetLPState(ctx, params.ConnectionId, params.IcaId, types.LPStateFunding)
    }

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
    // Safety net: if ICA address is already recorded, consider channel usable for FUNDING state
    if addr, ok := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, params.ConnectionId, portID); ok && addr != "" {
        if k.GetLPState(ctx, params.ConnectionId, params.IcaId) == types.LPStateUninitialized {
            k.SetLPState(ctx, params.ConnectionId, params.IcaId, types.LPStateFunding)
            k.Logger(ctx).Info("ICA address present; marking FUNDING", "port_id", portID, "addr", addr)
        }
    }
    // If we have a pending handshake recorded, wait for ack instead of opening new channels.
    if ctx.KVStore(k.storeKey).Get(pendingChanKey(owner)) != nil {
        k.Logger(ctx).Debug("EnsureICA: pending handshake exists, waiting for ack", "owner", owner)
        return
    }
    if _, found := k.icaControllerKeeper.GetActiveChannelID(ctx, params.ConnectionId, portID); found {
        // Mark state as FUNDING if not yet advanced
        if k.GetLPState(ctx, params.ConnectionId, params.IcaId) == types.LPStateUninitialized {
            k.SetLPState(ctx, params.ConnectionId, params.IcaId, types.LPStateFunding)
        }
        k.Logger(ctx).Info("ICA channel already active", "connection_id", params.ConnectionId, "ica_id", params.IcaId, "port_id", portID)
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

// DriveLPState performs lightweight state transitions using ICQ balances.
// - Registers ICQ balance queries for required host denoms when in FUNDING.
// - When balances meet thresholds, marks READY (pool creation handled later).
func (k Keeper) DriveLPState(ctx sdk.Context) {
    params := k.GetParams(ctx)
    if !params.AutoCreatePoolEnabled || params.ConnectionId == "" || params.IcaId == "" {
        return
    }
    // Throttle evaluation to align with ICQ update cadence
    if params.IcqUpdatePeriod > 0 {
        if uint64(ctx.BlockHeight())%params.IcqUpdatePeriod != 0 {
            return
        }
    }
    state := k.GetLPState(ctx, params.ConnectionId, params.IcaId)
    if state != types.LPStateFunding && state != types.LPStateReady {
        return
    }
    k.Logger(ctx).Info("DriveLPState evaluating", "height", ctx.BlockHeight(), "state", string(k.GetLPState(ctx, params.ConnectionId, params.IcaId)))
    // Ensure ICQ queries exist for each required denom
    if err := k.ensureBalanceICQs(ctx, params); err != nil {
        k.Logger(ctx).Error("ensureBalanceICQs failed", "err", err)
        return
    }
    // Check thresholds
    if k.meetsCreatePoolThreshold(ctx, params) {
        k.MarkLPReady(ctx, params.ConnectionId, params.IcaId)
    }
}

// ensureBalanceICQs registers KV queries on host bank balances for required denoms.
func (k Keeper) ensureBalanceICQs(ctx sdk.Context, p types.Params) error {
    // resolve ICA address on host
    owner := k.ownerString(p.IcaId)
    portID, err := icatypes.NewControllerPortID(owner)
    if err != nil { return err }
    icaAddr, found := k.icaControllerKeeper.GetInterchainAccountAddress(ctx, p.ConnectionId, portID)
    if !found || icaAddr == "" { return fmt.Errorf("ica address not found for %s on %s", portID, p.ConnectionId) }

    // decode host address bytes using configured HRP
    bech32 := addrcodec.NewBech32Codec(p.HostAccountHrp)
    addrBz, err := bech32.StringToBytes(icaAddr)
    if err != nil { return fmt.Errorf("decode host addr: %w", err) }

    // Build bank store key prefix: 0x02 | addrLen | addr
    prefix := append([]byte{0x02}, addrtypes.MustLengthPrefix(addrBz)...)

    k.Logger(ctx).Info("Ensuring ICQ balance queries", "connection_id", p.ConnectionId, "ica_id", p.IcaId, "host_addr", icaAddr)

    for _, need := range p.CreatePoolRequired {
        hostDenom := need.Denom
        if p.BaseLocalDenom != "" && need.Denom == p.BaseLocalDenom {
            // Derive host denom if possible
            if p.BaseHostDenom != "" {
                hostDenom = p.BaseHostDenom
            } else if p.BaseDexChannel != "" {
                // transfer/<dex_channel>/<base_local_denom>
                prefixed := transfertypes.GetPrefixedDenom("transfer", p.BaseDexChannel, p.BaseLocalDenom)
                hostDenom = transfertypes.ParseDenomTrace(prefixed).IBCDenom()
                k.Logger(ctx).Info("Derived host denom for base_local_denom", "local", p.BaseLocalDenom, "dex_channel", p.BaseDexChannel, "host", hostDenom)
            }
        }
        if p.QuoteLocalDenom != "" && need.Denom == p.QuoteLocalDenom {
            if p.QuoteHostDenom != "" {
                hostDenom = p.QuoteHostDenom
            } else if p.QuoteDexChannel != "" {
                prefixed := transfertypes.GetPrefixedDenom("transfer", p.QuoteDexChannel, p.QuoteLocalDenom)
                hostDenom = transfertypes.ParseDenomTrace(prefixed).IBCDenom()
                k.Logger(ctx).Info("Derived host denom for quote_local_denom", "local", p.QuoteLocalDenom, "dex_channel", p.QuoteDexChannel, "host", hostDenom)
            }
        }
        if qid, ok := k.getBalanceQueryID(ctx, p.ConnectionId, p.IcaId, hostDenom); ok {
            k.Logger(ctx).Info("ICQ already registered for denom", "denom", hostDenom, "query_id", qid)
            continue
        }
        key := append(prefix, []byte(hostDenom)...)
        rq := &iqtypes.RegisteredQuery{
            Id:                 0, // will assign below
            Owner:              authtypes.NewModuleAddress(types.ModuleName).String(),
            QueryType:          string(iqtypes.InterchainQueryTypeKV),
            // Use standard SDK proof path expected by ICQ relayers
            Keys:               []*iqtypes.KVKey{{Path: "store/bank/key", Key: key}},
            TransactionsFilter: "",
            ConnectionId:       p.ConnectionId,
            UpdatePeriod:       p.IcqUpdatePeriod,
            Deposit:            nil,
            SubmitTimeout:      k.icqKeeper.GetParams(ctx).QuerySubmitTimeout,
            RegisteredAtHeight: uint64(ctx.BlockHeight()),
        }
        last := k.icqKeeper.GetLastRegisteredQueryKey(ctx)
        last++
        k.icqKeeper.SetLastRegisteredQueryKey(ctx, last)
        rq.Id = last
        if err := k.icqKeeper.SaveQuery(ctx, rq); err != nil {
            return err
        }
        k.setBalanceQueryID(ctx, p.ConnectionId, p.IcaId, hostDenom, rq.Id)
        k.Logger(ctx).Info("Registered ICQ balance query", "denom", hostDenom, "query_id", rq.Id, "update_period", p.IcqUpdatePeriod)
    }
    return nil
}

// meetsCreatePoolThreshold returns true if all required host denoms meet thresholds.
func (k Keeper) meetsCreatePoolThreshold(ctx sdk.Context, p types.Params) bool {
    allMet := true
    for _, need := range p.CreatePoolRequired {
        hostDenom := need.Denom
        if p.BaseLocalDenom != "" && need.Denom == p.BaseLocalDenom {
            if p.BaseHostDenom != "" {
                hostDenom = p.BaseHostDenom
            } else if p.BaseDexChannel != "" {
                prefixed := transfertypes.GetPrefixedDenom("transfer", p.BaseDexChannel, p.BaseLocalDenom)
                hostDenom = transfertypes.ParseDenomTrace(prefixed).IBCDenom()
            }
        }
        if p.QuoteLocalDenom != "" && need.Denom == p.QuoteLocalDenom {
            if p.QuoteHostDenom != "" {
                hostDenom = p.QuoteHostDenom
            } else if p.QuoteDexChannel != "" {
                prefixed := transfertypes.GetPrefixedDenom("transfer", p.QuoteDexChannel, p.QuoteLocalDenom)
                hostDenom = transfertypes.ParseDenomTrace(prefixed).IBCDenom()
            }
        }
        qid, ok := k.getBalanceQueryID(ctx, p.ConnectionId, p.IcaId, hostDenom)
        if !ok {
            k.Logger(ctx).Info("ICQ query id not found yet", "denom", hostDenom)
            return false
        }
        res, err := k.icqKeeper.GetQueryResultByID(ctx, qid)
        if err != nil || res == nil || len(res.GetKvResults()) == 0 {
            k.Logger(ctx).Info("ICQ result not available yet", "denom", hostDenom, "query_id", qid)
            return false
        }
        val := res.GetKvResults()[0].GetValue()
        // Try to decode balance as sdk.Int (current), fallback to sdk.Coin (legacy)
        var amt sdkmath.Int
        if err := amt.Unmarshal(val); err != nil {
            var c sdk.Coin
            if e2 := c.Unmarshal(val); e2 == nil {
                amt = c.Amount
            } else {
                k.Logger(ctx).Error("failed to decode balance value", "denom", need.Denom, "err", err, "err2", e2)
                return false
            }
        }
        k.Logger(ctx).Info("ICQ balance", "denom", hostDenom, "current", amt.String(), "required", need.Amount.String())
        if amt.LT(need.Amount) { allMet = false }
    }
    if allMet {
        k.Logger(ctx).Info("All create_pool_required thresholds met; marking READY")
    }
    return allMet
}
