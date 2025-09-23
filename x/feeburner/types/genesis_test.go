package types_test

import (
	"testing"

	"github.com/maany-xyz/maany-app/app/config"

	"github.com/stretchr/testify/require"

	"github.com/maany-xyz/maany-app/testutil/common/nullify"
	keepertest "github.com/maany-xyz/maany-app/testutil/feeburner/keeper"
	"github.com/maany-xyz/maany-app/x/feeburner"
	"github.com/maany-xyz/maany-app/x/feeburner/types"
)

func TestGenesis(t *testing.T) {
	_ = config.GetDefaultConfig()

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	k, ctx := keepertest.FeeburnerKeeper(t)
	feeburner.InitGenesis(ctx, *k, genesisState)
	got := feeburner.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)
}

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "empty neutron denom",
			genState: &types.GenesisState{
				Params: types.Params{
					NeutronDenom: "",
				},
			},
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
