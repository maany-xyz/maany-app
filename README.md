# Maany App

- Binary: `maanyappd`
- Stack: Cosmos SDK, IBC (ICS), CosmWasm

This repo has been refactored to remove Osmosis DEX modules and custom MintBurn/GenesisMint modules. The remaining app wires core SDK + IBC, CosmWasm, and Neutron-derived utility modules (interchaintxs, interchainqueries, ibc-hooks, ibc-rate-limit, globalfee, feerefunder, feeburner, cron, contractmanager) and a wrapped IBC transfer module.

Build
- make build or make install (installs `maanyappd`)

Run localnet
- ./network/init.sh && ./network/init-neutrond.sh && ./network/start.sh
