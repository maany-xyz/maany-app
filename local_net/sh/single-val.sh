#!/bin/bash

set -euo pipefail

cd "$(dirname "$0")"

source ./utils.sh
# Resolve maanyappd binary: prefer local build, then PATH, then fallback
BINARY="${BINARY:-}"
if [ -z "${BINARY}" ]; then
  if [ -x "../build/maanyappd" ]; then
    BINARY="../build/maanyappd"
  elif command -v maanyappd >/dev/null 2>&1; then
    BINARY="maanyappd"
  elif [ -x "../maanyappd" ]; then
    BINARY="../maanyappd"
  else
    echo "maanyappd binary not found. Build with 'make build' first." >&2
    exit 1
  fi
fi

# Create local config directory inside local_net
HOME1=./config
echo "${HOME1}"
$BINARY -h

CCV_MODE=true
for arg in "$@"; do
  case $arg in
    -r|--reset)
      rm -rf "$HOME1"
      shift
      ;;
    --ccv)
      CCV_MODE=true
      shift
      ;;
  esac
done


echo "continuing"
if [ ! -f "$HOME1/data/priv_validator_state.json" ]; then

  "$BINARY" init validator --chain-id "maanyapp-local-1" --home "$HOME1" >/dev/null 2>&1
  # Apply minimal genesis adjustments
  update_genesis "$HOME1"

  if [ "$CCV_MODE" = false ]; then
    # Standalone single-validator mode (no CCV): create local validator and gentx
    "$BINARY" keys add validator --home "$HOME1" --keyring-backend test >/dev/null 2>&1
    "$BINARY" genesis add-genesis-account validator 1000000000uapp --home "$HOME1" --keyring-backend test
    "$BINARY" genesis gentx validator 100000000uapp --chain-id "local-maany-1" --home "$HOME1" --keyring-backend test >/dev/null 2>&1
    "$BINARY" genesis collect-gentxs --home "$HOME1" >/dev/null 2>&1
  else
    echo "Initialized consumer (CCV) genesis without local validators."
    echo "Ensure provider chain setup and CCV channel/proposal to receive validator set."
  fi
fi

echo "Local config initialized at $(cd "$HOME1" && pwd)"
echo "To start the node:"
echo "  $BINARY start --home $HOME1"
if [ "$CCV_MODE" = true ]; then
  echo "(CCV mode) After starting, connect this consumer to the provider via CCV:"
  echo "  - Start provider chain and relayer"
  echo "  - Submit Consumer addition on provider; establish CCV channel"
  echo "  - Consumer will adopt provider validator set"
fi
