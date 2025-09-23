# DEX COMMANDS

maanyappd tx bank send <from> <to> 1000untrn --keyring-backend test --gas auto --fees 500untrn

APP_HASH_B64=$(curl -s "$NODE/block?height=7" | jq -r '.result.block.header.app_hash')
APP_HASH_B64=$(echo "$APP_HASH_HEX" | xxd -r -p | base64 | tr -d '\n')
echo "$APP_HASH_B64"
