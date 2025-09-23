package keeper_test

import (
	"testing"

	"github.com/maany-xyz/maany-app/app/config"

	"github.com/maany-xyz/maany-app/testutil"

	testkeeper "github.com/maany-xyz/maany-app/testutil/cron/keeper"

	"github.com/stretchr/testify/require"

	"github.com/maany-xyz/maany-app/x/cron/types"
)

func TestGetParams(t *testing.T) {
	_ = config.GetDefaultConfig()

	k, ctx := testkeeper.CronKeeper(t, nil, nil)
	params := types.Params{
		SecurityAddress: testutil.TestOwnerAddress,
		Limit:           5,
	}

	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	require.EqualValues(t, params, k.GetParams(ctx))
}
