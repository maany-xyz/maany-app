package types

import (
    "context"

    wraptypes "github.com/maany-xyz/maany-app/x/transfer/types"
    ictxtypes "github.com/maany-xyz/maany-app/x/interchaintxs/types"
)

// TransferKeeper exposes the wrapped transfer MsgServer method
type TransferKeeper interface {
    Transfer(ctx context.Context, msg *wraptypes.MsgTransfer) (*wraptypes.MsgTransferResponse, error)
}

// InterchainTxsQueryKeeper exposes ICA address resolver
type InterchainTxsQueryKeeper interface {
    InterchainAccountAddress(c context.Context, req *ictxtypes.QueryInterchainAccountAddressRequest) (*ictxtypes.QueryInterchainAccountAddressResponse, error)
}

