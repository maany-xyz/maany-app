package keeper

import (
    "context"
    "fmt"
    "time"

    sdk "github.com/cosmos/cosmos-sdk/types"
    clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
    ictxtypes "github.com/maany-xyz/maany-app/x/interchaintxs/types"
    wraptypes "github.com/maany-xyz/maany-app/x/transfer/types"
    "github.com/maany-xyz/maany-app/x/autolp/types"
)

type msgServer struct{ Keeper }

var _ types.MsgServer = msgServer{}

func (m msgServer) CreateAutoLP(goCtx context.Context, msg *types.MsgCreateAutoLP) (*types.MsgCreateAutoLPResponse, error) {
    // manual validation until ValidateBasic is generated/customized
    if _, err := sdk.AccAddressFromBech32(msg.FromAddress); err != nil {
        return nil, err
    }
    if msg.ConnectionId == "" || msg.InterchainAccountId == "" || msg.TransferChannel == "" {
        return nil, fmt.Errorf("connection_id, interchain_account_id and transfer_channel must be set")
    }
    if !msg.Amount.IsValid() || msg.Amount.IsZero() {
        return nil, fmt.Errorf("invalid amount: %s", msg.Amount)
    }
    // Resolve ICA address from interchaintxs keeper
    icaResp, err := m.interchainQuery.InterchainAccountAddress(goCtx, &ictxtypes.QueryInterchainAccountAddressRequest{
        OwnerAddress:        msg.FromAddress,
        ConnectionId:        msg.ConnectionId,
        InterchainAccountId: msg.InterchainAccountId,
    })
    if err != nil {
        return nil, err
    }

    // Compose ICS-20 transfer to ICA address
    // Timeouts: translate seconds to timestamp if height timeout is zero
    timeoutTs := uint64(time.Now().Add(time.Duration(msg.TimeoutSeconds) * time.Second).UnixNano())
    // SourcePort is always transfer
    transferMsg := &wraptypes.MsgTransfer{
        SourcePort:       "transfer",
        SourceChannel:    msg.TransferChannel,
        Token:            msg.Amount,
        Sender:           msg.FromAddress,
        Receiver:         icaResp.InterchainAccountAddress,
        TimeoutHeight:    clienttypes.Height{},
        TimeoutTimestamp: timeoutTs,
        Memo:             "autolp:fund-ica",
    }

    resp, err := m.transferKeeper.Transfer(goCtx, transferMsg)
    if err != nil {
        return nil, err
    }

    return &types.MsgCreateAutoLPResponse{
        IcaAddress: icaResp.InterchainAccountAddress,
        SequenceId: resp.SequenceId,
        Channel:    resp.Channel,
    }, nil
}
