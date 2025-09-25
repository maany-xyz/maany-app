package main

import (
    "os"

    servercmd "github.com/cosmos/cosmos-sdk/server/cmd"
    "github.com/cosmos/cosmos-sdk/version"

    app "github.com/maany-xyz/maany-app/app"
)

func main() {
    rootCmd, _ := NewRootCmd()
    if err := servercmd.Execute(rootCmd, version.AppName, app.DefaultNodeHome); err != nil {
        os.Exit(1)
    }
}
