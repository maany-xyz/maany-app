# autolp Module components & testing

1. IBC hooks for autolp

- x/autolp/ibc_module.go
  - NewIBCModule(keeper.Keeper) implements porttypes.IBCModule with:
    - OnChanOpenAck → keeper.HandleChanOpenAck
    - OnAcknowledgementPacket → keeper.HandleAcknowledgement
    - OnTimeoutPacket → keeper.HandleTimeout
    - OnRecvPacket returns an error acknowledgement (no inbound ICA packets expected here).
  - CombineIBCModules(a, b) returns an IBCModule that forwards to both modules in order; we use it to invoke both interchaintxs’ and autolp’s handlers.

2. Keeper event emitters

- x/autolp/keeper/keeper.go
  - HandleChanOpenAck: emits autolp_ica_chan_open_ack event.
  - HandleAcknowledgement: emits autolp_ica_ack event.
  - HandleTimeout: emits autolp_ica_timeout event.

3. App wiring (router)

- app/app.go
  - Built the controller stack by combining interchaintxs and autolp IBC modules, then wrapping with icacontroller middleware:
    - base := interchaintxs.NewIBCModule(app.InterchainTxsKeeper)
    - autolpIBC := autolp.NewIBCModule(app.AutolpKeeper)
    - combined := autolp.CombineIBCModules(base, autolpIBC)
    - icaControllerStack := icacontroller.NewIBCMiddleware(combined, app.ICAControllerKeeper)
  - Router still routes icacontroller/ interchaintxs ports to icaControllerStack.

How the process works

- Register ICA (autolp.MsgRegisterICA, authority-gated):
  - Opens the ICS‑27 channel via ICA controller; relayer completes handshake.
  - autolp IBC module gets OnChanOpenAck and emits “autolp_ica_chan_open_ack”.
- Resolve ICA address:
  - autolp.QueryInterchainAccountAddress returns the host chain address (using ICA controller keeper).
- Fund ICA (autolp.MsgCreateAutoLP):
  - Sends ICS‑20 transfer over the transfer channel to the ICA address (DEX side).
- Submit ICA tx (autolp.MsgSubmitICATx, authority-gated):
  - Takes your host-chain messages as protobuf Any, serializes CosmosTx, calls ICA controller SendTx.
  - On ack/timeout, autolp emits “autolp_ica_ack” or “autolp_ica_timeout”.

Testing the logic

Prereqs

- IBC client+connection between your chain and the DEX.
- Relayer running for ICS‑27 (ICA) and ICS‑20.
- DEX icahost allow_messages includes your target type URLs (e.g., Osmosis GAMM).

Steps

1. Register the ICA

- Submit autolp.MsgRegisterICA via governance/admin:
  - { authority, connection_id, interchain_account_id, ordering: ORDER_ORDERED }
- Observe events on your chain:
  - autolp_ica_chan_open_ack with channel details.
- Check relayer logs for the handshake.

2. Resolve the ICA address

- Query autolp InterchainAccountAddress:
  - { connection_id, interchain_account_id }
- Save the returned ICA address (host chain).

3. Fund the ICA (ICS‑20)

- Submit autolp.MsgCreateAutoLP:
  - from_address (user), connection_id, interchain_account_id, transfer_channel, amount, timeout_seconds.
- Confirm the ICA address holds the IBC denom on the DEX.

4. Submit ICA tx (e.g., create a pool)

- Outside this repo, build the Any-encoded host message(s) using the DEX protos:
  - E.g., Any{type_url: “/osmosis.gamm.poolmodels.balancer.v1beta1.MsgCreateBalancerPool”, value: base64 bytes}, Sender set to the ICA address.
- Submit autolp.MsgSubmitICATx:
  - { authority, connection_id, interchain_account_id, msgs: [Any...], timeout_seconds }
- Observe your chain’s events:
  - On ack: autolp_ica_ack with port/channel/sequence/relayer
  - On timeout: autolp_ica_timeout
- Inspect the DEX chain state for pool creation or execution result.

Event reference

- autolp_ica_chan_open_ack
  - port_id, channel_id, cp_channel_id, cp_version
- autolp_ica_ack
  - src_port, src_channel, sequence, relayer
- autolp_ica_timeout
  - src_port, src_channel, sequence, relayer

Client reference (building Any externally)

- Use Osmosis protos in your client tool to build a concrete message, set Sender to the ICA address, then:
  - any, \_ := codectypes.NewAnyWithValue(msg)
  - Submit: autolptypes.MsgSubmitICATx{ Authority, ConnectionId, InterchainAccountId, Msgs: []\*codectypes.Any{any}, TimeoutSeconds }

Notes

- We intentionally did not vendor Osmosis protos into this repo (Option A). Clients construct Anys elsewhere and autolp forwards them.
- If you want a CLI helper inside this repo later, we can add Osmosis protos and build messages here.

Want me to also add:

- A tiny Go test/example that submits a dummy Any to verify ack/timeout events (using a mock host or a real devnet)?
- An optional key-value store in autolp to persist last-seen ack/timeout status per packet (for easy queries)?
