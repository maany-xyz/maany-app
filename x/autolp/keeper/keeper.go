package keeper

import (
    "context"
    "fmt"

    "cosmossdk.io/log"
    storetypes "cosmossdk.io/store/types"
    "github.com/cosmos/cosmos-sdk/codec"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/maany-xyz/maany-app/x/autolp/types"
    icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
    icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
    icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
    channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"
)

type Keeper struct {
    cdc      codec.BinaryCodec
    storeKey storetypes.StoreKey

    transferKeeper    types.TransferKeeper
    interchainQuery   types.InterchainTxsQueryKeeper

    // ICA controller
    icaControllerKeeper    icacontrollerkeeper.Keeper
    icaControllerMsgServer icacontrollertypes.MsgServer

    // authority address (e.g., gov/admin module address)
    authority string
}

func NewKeeper(
    cdc codec.BinaryCodec,
    key storetypes.StoreKey,
    transferKeeper types.TransferKeeper,
    interchainQuery types.InterchainTxsQueryKeeper,
    icaControllerKeeper icacontrollerkeeper.Keeper,
    icaControllerMsgServer icacontrollertypes.MsgServer,
    authority string,
) Keeper {
    return Keeper{
        cdc:            cdc,
        storeKey:       key,
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

// params: currently no module parameters are stored; return defaults
func (k Keeper) GetParams(ctx sdk.Context) types.Params { return types.Params{} }

func (k Keeper) SetParams(ctx sdk.Context, p types.Params) {}

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

func (k Keeper) ownerString(icaID string) string { return fmt.Sprintf("autolp/%s", icaID) }

// IBC event handlers (minimal): emit events for acks/timeouts/open-ack
func (k Keeper) HandleChanOpenAck(ctx sdk.Context, portID, channelID, counterpartyChannelID, counterpartyVersion string) error {
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
    ctx.EventManager().EmitEvent(sdk.NewEvent(
        "autolp_ica_timeout",
        sdk.NewAttribute("src_port", packet.SourcePort),
        sdk.NewAttribute("src_channel", packet.SourceChannel),
        sdk.NewAttribute("sequence", fmt.Sprintf("%d", packet.Sequence)),
        sdk.NewAttribute("relayer", relayer.String()),
    ))
    return nil
}
