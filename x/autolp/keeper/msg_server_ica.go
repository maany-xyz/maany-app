package keeper

import (
    "context"
    "time"

    "cosmossdk.io/errors"
    sdk "github.com/cosmos/cosmos-sdk/types"
    sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
    icacontrollertypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/types"
    icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
    channeltypes "github.com/cosmos/ibc-go/v8/modules/core/04-channel/types"

    "github.com/maany-xyz/maany-app/x/autolp/types"
)

// NewMsgServerImpl autolp MsgServer
func NewMsgServerImpl(k Keeper) types.MsgServer { return msgServer{k} }

// RegisterICA opens an ICA channel by registering an interchain account.
func (m msgServer) RegisterICA(goCtx context.Context, msg *types.MsgRegisterICA) (*types.MsgRegisterICAResponse, error) {
    if !m.isAuthorized(goCtx, msg.Authority) {
        return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s", m.authority)
    }
    ctx := sdk.UnwrapSDKContext(goCtx)
    owner := m.ownerString(msg.InterchainAccountId)

    // empty version uses default, ordering passed through
    resp, err := m.icaControllerMsgServer.RegisterInterchainAccount(ctx, &icacontrollertypes.MsgRegisterInterchainAccount{
        Owner:        owner,
        ConnectionId: msg.ConnectionId,
        Version:      "",
        Ordering:     channeltypes.Order(msg.Ordering),
    })
    if err != nil {
        return nil, err
    }
    return &types.MsgRegisterICAResponse{ChannelId: resp.ChannelId, PortId: resp.PortId}, nil
}

// SubmitICATx executes msgs on the host chain from the ICA.
func (m msgServer) SubmitICATx(goCtx context.Context, msg *types.MsgSubmitICATx) (*types.MsgSubmitICATxResponse, error) {
    if !m.isAuthorized(goCtx, msg.Authority) {
        return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s", m.authority)
    }
    ctx := sdk.UnwrapSDKContext(goCtx)
    owner := m.ownerString(msg.InterchainAccountId)

    // validate Any list is not empty
    if len(msg.Msgs) == 0 {
        return nil, errors.Wrapf(sdkerrors.ErrInvalidRequest, "no msgs provided")
    }

    // Serialize CosmosTx
    cosmosTx := &icatypes.CosmosTx{Messages: msg.Msgs}
    bz, err := m.cdc.Marshal(cosmosTx)
    if err != nil {
        return nil, err
    }

    packet := icatypes.InterchainAccountPacketData{
        Type: icatypes.EXECUTE_TX,
        Data: bz,
        Memo: msg.Memo,
    }

    resp, err := m.icaControllerMsgServer.SendTx(ctx, &icacontrollertypes.MsgSendTx{
        Owner:           owner,
        ConnectionId:    msg.ConnectionId,
        PacketData:      packet,
        RelativeTimeout: uint64(time.Duration(msg.TimeoutSeconds) * time.Second),
    })
    if err != nil {
        return nil, err
    }

    // For channel, try to resolve active channel ID
    portID, _ := icatypes.NewControllerPortID(owner)
    channelID, _ := m.icaControllerKeeper.GetActiveChannelID(ctx, msg.ConnectionId, portID)

    return &types.MsgSubmitICATxResponse{SequenceId: resp.Sequence, Channel: channelID}, nil
}

// isAuthorized returns true if addr equals module authority or is on the allowlist.
func (m msgServer) isAuthorized(goCtx context.Context, addr string) bool {
    if addr == m.authority { return true }
    ctx := sdk.UnwrapSDKContext(goCtx)
    params := m.GetParams(ctx)
    for _, a := range params.AllowedSubmitters {
        if a == addr { return true }
    }
    return false
}

// UpdateParams updates module parameters; allowed for module authority or any
// address present in the current allowlist.
func (m msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
    if !m.isAuthorized(goCtx, msg.Authority) {
        return nil, errors.Wrapf(sdkerrors.ErrUnauthorized, "invalid authority; expected %s or allowed submitter", m.authority)
    }
    if err := msg.Params.Validate(); err != nil {
        return nil, err
    }
    ctx := sdk.UnwrapSDKContext(goCtx)
    m.SetParams(ctx, msg.Params)
    return &types.MsgUpdateParamsResponse{}, nil
}
