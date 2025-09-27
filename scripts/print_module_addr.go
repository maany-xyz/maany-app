// print_module_addr.go
// Tiny helper to derive a module account address (bech32) offline.
// Usage:
//   go run ./scripts/print_module_addr.go \
//     -module autolp \
//     -prefix maany \
//     -pubprefix maanypub
// Notes:
// - The bech32 account prefix must match your chain config.
// - The pubkey prefix is typically <prefix> + "pub".

package main

import (
	"flag"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func main() {
    var (
        module    = flag.String("module", "autolp", "module account name (e.g. autolp)")
        prefix    = flag.String("prefix", "maanyapp", "bech32 account address prefix")
        pubprefix = flag.String("pubprefix", "maanyapppub", "bech32 account pubkey prefix")
    )
    flag.Parse()

    cfg := sdk.GetConfig()
    cfg.SetBech32PrefixForAccount(*prefix, *pubprefix)
    // Optional: if you also need validator/consensus prefixes, set them too.

    addr := authtypes.NewModuleAddress(*module)
    fmt.Println(addr.String())
}

