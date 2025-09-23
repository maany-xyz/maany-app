#!/bin/bash

# Minimal genesis adjustments for this project
# - Set staking bond denom to uapp
# - Optionally add bank metadata for uapp
update_genesis() {
  local HOME1=$1

  local GEN=$HOME1/config/genesis.json
  local TEMP="$HOME1/genesis.tmp.json"

  # Ensure temp exists
  : > "$TEMP"

  # Set staking bond denom to uapp
  # jq '.app_state.staking.params.bond_denom = "uapp"' "$GEN" > "$TEMP" && mv "$TEMP" "$GEN"

  # Ensure bank metadata for uapp exists
  jq '.app_state.bank.denom_metadata += [{ "description": "Maany base denom", "denom_units": [{ "denom": "uapp", "exponent": 0, "aliases": ["microuapp"] }, { "denom": "APP", "exponent": 6 }], "base": "uapp", "display": "APP", "name": "Maany", "symbol": "APP" }]' "$GEN" > "$TEMP" && mv "$TEMP" "$GEN"
}
