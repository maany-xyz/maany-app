package main

import (
    "fmt"
    "os"
)

// This is a placeholder binary to keep the repo buildable.
// The full CLI (init/start/export, etc.) should be restored from your original cmd tree
// or scaffolded to wire Cosmos SDK server commands to the App in app/.
func main() {
    fmt.Fprintln(os.Stderr, "maanyappd CLI is not wired yet in this repository.\n- If you have the original cmd/ contents, remove any .gitignore rules that exclude it and add them back.\n- Otherwise, I can scaffold a minimal Cosmos SDK root command on request.")
    os.Exit(1)
}

