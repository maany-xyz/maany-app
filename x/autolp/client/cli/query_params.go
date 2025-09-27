package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/spf13/cobra"

    "github.com/maany-xyz/maany-app/x/autolp/types"
)

// CmdQueryParams returns the command to query autolp module params.
func CmdQueryParams() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "params",
        Short: "Get the params for the x/autolp module",
        Args:  cobra.NoArgs,
        RunE: func(cmd *cobra.Command, _ []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }

            queryClient := types.NewQueryClient(clientCtx)
            res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
            if err != nil {
                return err
            }
            return clientCtx.PrintProto(res)
        },
    }
    flags.AddQueryFlagsToCmd(cmd)
    return cmd
}

