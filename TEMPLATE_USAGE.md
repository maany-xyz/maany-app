Template usage

- Requirements: cookiecutter (and optionally cue)

1) Generate a new app repository

   cookiecutter https://github.com/maany-xyz/maany-app

   Parameters:
   - project_name: Display name of the app
   - module_path: New Go module path (e.g., github.com/you/foo)
   - binary_name: CLI binary (e.g., fooappd)
   - bech32_main_prefix: Address prefix (e.g., foo)
   - base_denom: Base denom (e.g., ufoo)
   - display_denom: Display denom (e.g., FOO)
   - denom_exponent: Display exponent (e.g., 6)
   - chain_id: Default chain-id for local scripts
   - min_gas_price: Default min gas price for templates
   - home_dir_name: Home directory (dot-prefixed)

2) Post-gen steps

- The post-gen hook will rewrite import paths from the original
  module root to your chosen module_path and set go.mod accordingly.
- Run make build or make install in the generated project.

3) Validation (optional)

- If you have cue installed, you can validate the chosen parameters:

  bash template/validate.sh

Notes

- Chain ID is not compiled; scripts will use the provided default.
- Some docs and tests may still reference sample values; update as needed.

