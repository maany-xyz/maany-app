#!/usr/bin/env bash
set -euo pipefail

# Renders node config templates (genesis.json, app.toml, config.toml)
# Requires gomplate installed.

if ! command -v gomplate >/dev/null 2>&1; then
  echo "gomplate not found; install from https://gomplate.ca/" >&2
  exit 1
fi

OUT_DIR=${1:-./rendered}
mkdir -p "$OUT_DIR"

# Defaults from Cookiecutter values at generation time
export CHAIN_ID=${CHAIN_ID:-"{{ cookiecutter.chain_id }}"}
export BASE_DENOM=${BASE_DENOM:-"{{ cookiecutter.base_denom }}"}
export DISPLAY_DENOM=${DISPLAY_DENOM:-"{{ cookiecutter.display_denom }}"}
export DENOM_EXP=${DENOM_EXP:-"{{ cookiecutter.denom_exponent }}"}
export MIN_GAS_PRICE=${MIN_GAS_PRICE:-"{{ cookiecutter.min_gas_price }}"}

gomplate -f genesis.json.tmpl -o "$OUT_DIR/genesis.json"
gomplate -f app.toml.tmpl   -o "$OUT_DIR/app.toml"
gomplate -f config.toml.tmpl -o "$OUT_DIR/config.toml"

echo "Rendered configs to $OUT_DIR"

