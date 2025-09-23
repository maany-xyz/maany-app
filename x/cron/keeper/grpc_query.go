package keeper

import (
	"github.com/maany-xyz/maany-app/x/cron/types"
)

var _ types.QueryServer = Keeper{}
