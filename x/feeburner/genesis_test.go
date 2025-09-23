package feeburner_test

import (
	"testing"

	"cosmossdk.io/math"

	"github.com/maany-xyz/maany-app/app/config"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/maany-xyz/maany-app/testutil/common/nullify"
	"github.com/maany-xyz/maany-app/testutil/feeburner/keeper"
	"github.com/maany-xyz/maany-app/x/feeburner"
	"github.com/maany-xyz/maany-app/x/feeburner/types"
)

func TestGenesis(t *testing.T) {
	_ = config.GetDefaultConfig()

	amount := math.NewInt(10)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		TotalBurnedNeutronsAmount: types.TotalBurnedNeutronsAmount{
			Coin: sdk.NewCoin(types.DefaultNeutronDenom, amount),
		},
	}

	k, ctx := keeper.FeeburnerKeeper(t)
	feeburner.InitGenesis(ctx, *k, genesisState)

	burnedTokens := k.GetTotalBurnedNeutronsAmount(ctx)
	require.Equal(t, amount, burnedTokens.Coin.Amount)

	got := feeburner.ExportGenesis(ctx, *k)
	require.NotNil(t, got)
	require.Equal(t, amount, got.TotalBurnedNeutronsAmount.Coin.Amount)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}
