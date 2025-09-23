package keeper

import (
	"github.com/maany-xyz/maany-app/x/interchaintxs/types"
)

var _ types.QueryServer = Keeper{}
