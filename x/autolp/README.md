# autolp module

Overview

- autolp is a controller-side module that orchestrates liquidity provisioning on a remote DEX chain via Interchain Accounts (ICA) and funding via ICS-20 transfers.
- It does not execute DEX logic locally. Instead it packages host-chain messages (protobuf Any) into an ICA packet; the host (DEX) chain executes them (e.g., Osmosis GAMM messages).

High-level architecture

- Controller (this chain):
  - Register an ICA over an IBC connection to the DEX (ICS-27 channel).
  - Fund the ICA on the DEX via ICS-20 (transfer channel).
  - Submit host-chain messages (as Any) for the ICA to execute on the DEX.
- Host (DEX chain):
  - Must have icahost enabled and allow the target message type URLs.
  - Executes messages under the ICA address.

Messages (tx)

- MsgRegisterICA

  - authority, connection_id, interchain_account_id, ordering (ORDER_ORDERED recommended)
  - Opens the ICS-27 ICA channel; the relayer completes the handshake.

- MsgCreateAutoLP

  - from_address, connection_id, interchain_account_id, transfer_channel, amount (Coin), timeout_seconds
  - Sends an ICS-20 transfer from from_address to the ICA address on the DEX (funding step).

- MsgSubmitICATx

  - authority, connection_id, interchain_account_id, msgs ([]Any), memo, timeout_seconds
  - Serializes msgs into icatypes.CosmosTx and sends via ICA. The DEX executes them.
  - msgs are host-chain messages encoded as protobuf Any (type_url + value). Construction can be done outside this repo.

- MsgUpdateParams
  - authority, params (see Parameters)
  - Updates module parameters on-chain. Authorized by module authority or any address in the allowlist (see Authorization).

Queries

- Query/Params: returns current module parameters.
- Query/InterchainAccountAddress: resolves the host-chain ICA address for (connection_id, interchain_account_id).

Parameters

- Params.allowed_submitters: []string (bech32 account addresses)
  - Addresses in this allowlist are permitted to submit autolp ICA actions directly.
  - The module authority (gov/admin module account) is always permitted in addition to this list.

Authorization model

- autolp enforces an authority check for the following messages: MsgRegisterICA, MsgSubmitICATx, MsgUpdateParams.
- A message is authorized if:
  1. msg.authority equals the module authority (gov/admin module account), or
  2. msg.authority is present in Params.allowed_submitters.
- This allows either governance-based control, a direct operator address (e.g., multisig), or both.

IBC integration and events

- autolp registers a lightweight IBC module for ICS-27 callbacks and emits events:
  - autolp_ica_chan_open_ack: on channel open ack (port_id, channel_id, cp_channel_id, cp_version)
  - autolp_ica_ack: on packet acknowledgement (src_port, src_channel, sequence, relayer)
  - autolp_ica_timeout: on packet timeout (src_port, src_channel, sequence, relayer)
- Controller routing is composed with the existing interchaintxs IBC module and wrapped by the ICA controller middleware.

Typical flow (manual)

1. Ensure IBC client+connection exists to the DEX and a relayer is running.
2. Create an ICS-20 transfer channel between the chains (relayer-driven, one-time).
3. Register the ICA:
   - Submit MsgRegisterICA with authority set to the module authority or an address from allowed_submitters.
   - Relayer completes handshake. Watch autolp_ica_chan_open_ack.
4. Query the ICA address via Query/InterchainAccountAddress.
5. Fund the ICA via MsgCreateAutoLP (ICS-20).
6. Submit host-chain messages via MsgSubmitICATx (e.g., Osmosis GAMM pool create/join), passing msgs as protobuf Any.
   - The DEX executes messages; watch autolp_ica_ack or autolp_ica_timeout.

Constructing host-chain messages (Any)

- autolp forwards Any without decoding. Construct them in your client/tool using the host-chain protobufs.
- Example Any JSON (submitted via gRPC):
  {
  "type_url": "/osmosis.gamm.poolmodels.balancer.v1beta1.MsgCreateBalancerPool",
  "value": "BASE64_PROTO_BYTES"
  }
- Sender in host messages should be the ICA address on the DEX chain.

Parameters: initialization and updates

- Genesis (initial):
  - In genesis.json under app_state.autolp.params:
    {
    "allowed_submitters": [
    "yourprefix1admin...",
    "yourprefix1multisig..."
    ]
    }
- On-chain update (MsgUpdateParams):
  - Authorized by module authority or any address listed in the current allowed_submitters.
  - Example body:
    {
    "@type": "/maany.autolp.v1.MsgUpdateParams",
    "authority": "yourprefix1admin...",
    "params": {
    "allowed_submitters": [
    "yourprefix1admin...",
    "yourprefix1multisig..."
    ]
    }
    }

Host-chain requirements

- icahost enabled; Params.allow_messages must include the type URLs you plan to execute (e.g., Osmosis GAMM msgs).
- A relayer must be running to relay ICS-27 (ICA) and ICS-20 packets.

Security considerations

- Keep allowed_submitters minimal and prefer multisig.
- Ensure host-chain allowlists are appropriately configured; otherwise SubmitICATx will be rejected.
- Timeouts are configurable per message; set conservatively to account for relayer latency.

Troubleshooting

- No ICA channel: ensure RegisterICA was submitted and a relayer is active.
- No funds on DEX: verify ICS-20 transfer channel and MsgCreateAutoLP execution.
- Host message rejected: check host allow_messages and that Any was encoded with the correct type_url and bytes.
- Observe emitted events (autolp_ica_*) to follow packet lifecycle.

Samples

- Governance proposal (MsgUpdateParams)
  - File: update_autolp_params.json
  {
    "messages": [
      {
        "@type": "/maany.autolp.v1.MsgUpdateParams",
        "authority": "yourprefix1moduleorauthorizedaddr...",
        "params": {
          "allowed_submitters": [
            "yourprefix1operatoraddr...",
            "yourprefix1multisigaddr..."
          ]
        }
      }
    ],
    "metadata": "",
    "deposit": "1000000untrn",
    "title": "Update autolp params",
    "summary": "Set allowed_submitters for direct ICA actions"
  }

  - Submit + vote (example):
    - maanyappd tx gov submit-proposal update_autolp_params.json --from <gov-key> --gas auto --chain-id <chain>
    - maanyappd tx gov vote <proposal-id> yes --from <gov-key> --chain-id <chain>

- SubmitICATx with Any (JSON body)
  - Build host-chain Any in your client/tool. Example body to send via gRPC/CLI:
  {
    "@type": "/maany.autolp.v1.MsgSubmitICATx",
    "authority": "yourprefix1operatoraddr...",
    "connection_id": "connection-0",
    "interchain_account_id": "dex-lp-ica",
    "msgs": [
      {
        "type_url": "/osmosis.gamm.poolmodels.balancer.v1beta1.MsgCreateBalancerPool",
        "value": "BASE64_PROTO_BYTES"
      }
    ],
    "memo": "autolp:create-balancer-pool",
    "timeout_seconds": 300
  }

  - Notes:
    - Replace BASE64_PROTO_BYTES with the protobuf-encoded bytes of the host message (e.g., MsgCreateBalancerPool with Sender set to the ICA address on the DEX).
    - The signer of the tx must match the authority field and be authorized (module authority or in allowed_submitters).
