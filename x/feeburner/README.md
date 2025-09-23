Here’s a practical, working flow to create a gauge that rewards locked tokens on this
chain.

1. Pick a lock duration

- See supported durations: - maanyappd q incentives lockable-durations --home local_net/config --node http://
  localhost:26657 --chain-id local-maany-1
- Choose one (example: 24h).

2. Create a ByDuration (lock-based) gauge

- Syntax (we patched CLI to only allow ByDuration; poolId must be 0):
  - maanyappd tx incentives create-gauge uapp 1uapp 0 \
    --duration 24h \
    --perpetual \
    --start-time 0 \
    --home local_net/config \
    --node http://localhost:26657 \
    --chain-id local-maany-1 \
    --keyring-backend test \
    --gas auto --gas-adjustment 1.3 --fees 1000uapp -y
- Notes: - lockup_denom: uapp - reward: you can start with a small funding (e.g., 1uapp). Feeburner will top up each
  block if you’ve configured routing. - poolId: must be 0 (pools are disabled). - --perpetual: true means the gauge distributes every incentives epoch and never
  finishes; if you want a finite schedule, omit --perpetual and pass --epochs N instead.

3. Get the gauge ID

- maanyappd q incentives gauges --home local_net/config --node http://localhost:26657
  --chain-id local-maany-1 -o json | jq '.data[] | {id, denom: .distribute_to.denom,
  duration: .distribute_to.duration, perpetual: .is_perpetual}'
- Copy the id value.

4. Configure fee routing (so consumer fees fund the gauge)

- We added code in feeburner that can route all native consumer fees to a configured
  incentives gauge. This is wired in code (Keeper.SetIncentivesRouting) not via CLI.
- Two ways to set it: - In code during app startup: call SetIncentivesRouting(app.IncentivesKeeper,
  <gaugeID>, <fraction or full>) - We currently route ALL native fees if a gaugeID is set (no burn). If you prefer
  a fraction param, I can wire it to feeburner params. - Alternatively, temporarily fund the gauge manually while we add param-based routing: - maanyappd tx incentives add-to-gauge --gauge-id <ID> --amount 100000uapp ...
  usual flags

5. Lock tokens (users)

- Lock coins for the target duration:
  - maanyappd tx lockup lock-tokens --amount 100000uapp --duration 24h \
    --home local_net/config --node http://localhost:26657 \
    --chain-id local-maany-1 --keyring-backend test \
    --gas auto --fees 1000uapp -y
- Users can query their locks: - maanyappd q lockup account-locked-pastime-not-unlocking $(maanyappd keys show <name>
  -a --home local_net/config --keyring-backend test) 24h ...

6. Verify distribution

- After an incentives epoch ends, check: - maanyappd q incentives active-gauges --home local_net/config --node http://
  localhost:26657 --chain-id local-maany-1 - maanyappd q incentives module-to-distribute-coins ...
- And user balances:
  - maanyappd q bank balances <addr> ...

Quick recap

- Create a perpetual ByDuration gauge for uapp with a supported duration and poolId=0.
- Route consumer fees to that gauge (code-based for now) or fund it manually.
- Users lock uapp for that duration.
- Incentives distribute on epoch end.

If you want, I can add a small feeburner param (gauge_id + stake_all_native) and a CLI tx
to configure routing without changing code.
