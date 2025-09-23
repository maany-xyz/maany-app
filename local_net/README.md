# Local Chain Setup

- Change `rpc` to free port
- Change `grpc` to free port

`maanyappd tx bank send maanyapp17ykqz4ec2ygt2a9e2jle2egxyr5spnx4kxxpjf maanyapp1h0c0hra2sy47ejchahpk6szldger7qegu0ntr8 100000000uapp \
  --chain-id maanyapp-local-1 \
  --node http://localhost:26687 \
  --home local_net/.maanyappd \
  --keyring-backend test \
  --gas auto --gas-adjustment 1.3 \
  --fees 100000uapp -y`

`maanyappd q bank balances maanyapp1x69dz0c0emw8m2c6kp5v6c08kgjxmu30y2yjp5 --node http://localhost:26687 `

`maanyappd comet unsafe-reset-all --home local_net/config`
