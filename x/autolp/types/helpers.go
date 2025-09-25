package types

import (
    codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

// NewAny packs raw bytes with a typeURL into a codectypes.Any suitable for MsgSubmitICATx.
// This is useful when you don’t have the concrete Go type registered locally but know the
// host chain typeURL and its serialized bytes.
func NewAny(typeURL string, bz []byte) *codectypes.Any {
    return &codectypes.Any{TypeUrl: typeURL, Value: bz}
}

