package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/spf13/cobra"

    "github.com/maany-xyz/maany-app/x/autolp/types"
)

// CmdInterchainAccountAddress queries the ICA address for a given owner, ICA ID and connection.
func CmdInterchainAccountAddress() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "interchain-account-address [owner-address] [interchain-account-id] [connection-id]",
        Short: "Resolve the ICA address on the host for owner, ica-id and connection",
        Long: `Resolve the interchain account address stored by the controller.

Arguments:
  owner-address           Bech32 address of the controller owner (e.g. autolp module address)
  interchain-account-id   Local ICA identifier used by the module (e.g. dex-lp-ica)
  connection-id           IBC connection on the controller (e.g. connection-0)

Example:
  maanyappd q autolp interchain-account-address maanyapp1... dex-lp-ica connection-0`,
        Args: cobra.ExactArgs(3),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientQueryContext(cmd)
            if err != nil {
                return err
            }

            queryClient := types.NewQueryClient(clientCtx)
            res, err := queryClient.InterchainAccountAddress(cmd.Context(), &types.QueryICAAddressRequest{
                OwnerAddress:        args[0],
                InterchainAccountId: args[1],
                ConnectionId:        args[2],
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

