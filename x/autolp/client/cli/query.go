package cli

import (
    "fmt"

    "github.com/cosmos/cosmos-sdk/client"
    "github.com/spf13/cobra"

    "github.com/maany-xyz/maany-app/x/autolp/types"
)

// GetQueryCmd returns the cli query commands for the autolp module
func GetQueryCmd(_ string) *cobra.Command {
    // Group autolp queries under a subcommand
    cmd := &cobra.Command{
        Use:                        types.ModuleName,
        Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
        DisableFlagParsing:         true,
        SuggestionsMinimumDistance: 1,
        RunE:                       client.ValidateCmd,
    }

    cmd.AddCommand(CmdInterchainAccountAddress())
    cmd.AddCommand(CmdLPState())
    cmd.AddCommand(CmdQueryParams())

    return cmd
}
