package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	"github.com/maany-xyz/maany-app/x/incentives/types"
)

func (k Keeper) HandleCreateGaugeProposal(ctx sdk.Context, p *types.CreateGroupsProposal) error {
    // Pools/groups are disabled in this build.
    return fmt.Errorf("create groups proposal is unsupported: pools disabled")
}

func NewIncentivesProposalHandler(k Keeper) govtypesv1.Handler {
	return func(ctx sdk.Context, content govtypesv1.Content) error {
		switch c := content.(type) {
		case *types.CreateGroupsProposal:
			return k.HandleCreateGaugeProposal(ctx, c)

		default:
			return fmt.Errorf("unrecognized incentives proposal content type: %T", c)
		}
	}
}
