package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/spf13/cobra"

    "github.com/maany-xyz/maany-app/x/autolp/types"
)

// CmdLPState queries the current LP state for (connection_id, ica_id).
func CmdLPState() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "lp-state [interchain-account-id] [connection-id]",
        Short: "Show autolp lifecycle state (UNINITIALIZED|FUNDING|READY|CREATING_POOL|ACTIVE|PAUSED)",
        Args:  cobra.ExactArgs(2),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }
            queryClient := types.NewQueryClient(clientCtx)
            res, err := queryClient.LPState(cmd.Context(), &types.QueryLPStateRequest{
                InterchainAccountId: args[0],
                ConnectionId:        args[1],
            })
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(res)
        },
    }
    flags.AddQueryFlagsToCmd(cmd)
    return cmd
}

