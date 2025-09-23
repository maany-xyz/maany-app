package keeper_test

import (
	"testing"

	"github.com/maany-xyz/maany-app/app/config"

	"github.com/stretchr/testify/require"

	testkeeper "github.com/maany-xyz/maany-app/testutil/feeburner/keeper"
	"github.com/maany-xyz/maany-app/x/feeburner/types"
)

func TestGetParams(t *testing.T) {
	_ = config.GetDefaultConfig()

	k, ctx := testkeeper.FeeburnerKeeper(t)
	params := types.DefaultParams()

	err := k.SetParams(ctx, params)
	require.NoError(t, err)

	require.EqualValues(t, params, k.GetParams(ctx))
}
