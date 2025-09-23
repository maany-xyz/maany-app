package keeper_test

import (
    "testing"

    "cosmossdk.io/math"
    sdktypes "github.com/cosmos/cosmos-sdk/types"

    testkeeper "github.com/maany-xyz/maany-app/testutil/globalfee/keeper"

    "github.com/stretchr/testify/require"

    "github.com/maany-xyz/maany-app/x/globalfee/types"
    appparams "github.com/maany-xyz/maany-app/app/params"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.GlobalFeeKeeper(t)
	wctx := ctx
	params := types.DefaultParams()
    params.MinimumGasPrices = sdktypes.NewDecCoins(sdktypes.NewDecCoin(appparams.DefaultDenom, math.NewInt(1)))
	err := keeper.SetParams(ctx, params)
	require.NoError(t, err)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
