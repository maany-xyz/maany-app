package autolp

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/maany-xyz/maany-app/x/autolp/keeper"
    "github.com/maany-xyz/maany-app/x/autolp/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
    k.SetParams(ctx, genState.Params)
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
    gs := types.DefaultGenesis()
    gs.Params = k.GetParams(ctx)
    return gs
}

