#!/usr/bin/env python3
"""
Post-generation hook for Cookiecutter.

Tasks:
- Replace the module import root "github.com/maany-xyz/maany-app" with the chosen module_path.
- Update go.mod module line to the chosen module_path if still the default.
- Apply chosen parameters to source files without requiring placeholders in the template source.
"""

import os
import sys

ROOT = os.getcwd()

DEFAULT_MODULE = "github.com/maany-xyz/maany-app"
TARGET_MODULE = "{{ cookiecutter.module_path }}"
BINARY_NAME = "{{ cookiecutter.binary_name }}"
BECH32_PREFIX = "{{ cookiecutter.bech32_main_prefix }}"
HUMAN_DENOM = "{{ cookiecutter.display_denom }}"
BASE_DENOM = "{{ cookiecutter.base_denom }}"
DENOM_EXP = "{{ cookiecutter.denom_exponent }}"
CHAIN_ID = "{{ cookiecutter.chain_id }}"


def is_text_file(path: str) -> bool:
    try:
        with open(path, "rb") as f:
            chunk = f.read(4096)
        # Heuristic: if it has null bytes, treat as binary
        return b"\x00" not in chunk
    except Exception:
        return False


def replace_in_file(path: str, before: str, after: str) -> None:
    try:
        with open(path, "r", encoding="utf-8") as f:
            data = f.read()
    except UnicodeDecodeError:
        return
    if before not in data:
        return
    data = data.replace(before, after)
    with open(path, "w", encoding="utf-8") as f:
        f.write(data)


def walk_and_replace(root: str, before: str, after: str) -> None:
    for dirpath, dirnames, filenames in os.walk(root):
        # Skip .git and build dirs
        if ".git" in dirnames:
            dirnames.remove(".git")
        if "build" in dirnames:
            dirnames.remove("build")
        for name in filenames:
            path = os.path.join(dirpath, name)
            # Skip obvious binaries and images
            if not is_text_file(path):
                continue
            replace_in_file(path, before, after)


def main() -> int:
    # Replace module import root across the project
    if TARGET_MODULE and TARGET_MODULE != DEFAULT_MODULE:
        walk_and_replace(ROOT, DEFAULT_MODULE, TARGET_MODULE)

    # Ensure go.mod has correct module line
    gomod = os.path.join(ROOT, "go.mod")
    if os.path.exists(gomod):
        replace_in_file(gomod, f"module {DEFAULT_MODULE}", f"module {TARGET_MODULE}")

    # Apply core params in code
    replace_in_file(os.path.join(ROOT, "app", "app.go"),
                    'Name = "maanyappd"', f'Name = "{BINARY_NAME}"')

    cfg_path = os.path.join(ROOT, "app", "config", "config.go")
    replace_in_file(cfg_path, 'HumanCoinUnit     = "APP"', f'HumanCoinUnit     = "{HUMAN_DENOM}"')
    replace_in_file(cfg_path, 'BaseCoinUnit      = "uapp"', f'BaseCoinUnit      = "{BASE_DENOM}"')
    replace_in_file(cfg_path, 'DefaultBondDenom  = BaseCoinUnit', 'DefaultBondDenom  = BaseCoinUnit')
    replace_in_file(cfg_path, 'DefaultExponent   = 6', f'DefaultExponent   = {DENOM_EXP}')
    replace_in_file(cfg_path, 'Bech32MainPrefix = "maanyapp"', f'Bech32MainPrefix = "{BECH32_PREFIX}"')

    replace_in_file(os.path.join(ROOT, "app", "params", "denom.go"),
                    'const DefaultDenom = "uapp"', f'const DefaultDenom = "{BASE_DENOM}"')
    replace_in_file(os.path.join(ROOT, "app", "genesis.go"),
                    'var FeeDenom = "uapp"', f'var FeeDenom = "{BASE_DENOM}"')

    # Makefile tweaks: AppName, and local run env vars
    makefile = os.path.join(ROOT, "Makefile")
    replace_in_file(makefile,
                    'version.Name=maanyapp',
                    'version.Name={{ cookiecutter.project_name | replace(" ", "") | lower }}')
    replace_in_file(makefile,
                    'version.AppName=maanyappd',
                    f'version.AppName={BINARY_NAME}')
    replace_in_file(makefile,
                    'Starting up maanyappd alone...',
                    f'Starting up {BINARY_NAME} alone...')
    replace_in_file(makefile,
                    'export BINARY=maanyappd',
                    f'export BINARY={BINARY_NAME}')
    replace_in_file(makefile,
                    'CHAINID=test-1',
                    f'CHAINID={CHAIN_ID}')
    replace_in_file(makefile,
                    'STAKEDENOM=uapp',
                    f'STAKEDENOM={BASE_DENOM}')
    replace_in_file(makefile,
                    'Killing maanyappd',
                    f'Killing {BINARY_NAME}')
    replace_in_file(makefile,
                    'BINARY=maanyappd .',
                    f'BINARY={BINARY_NAME} .')

    # Dockerfile: AppName and output binary/entrypoint
    dockerfile = os.path.join(ROOT, "Dockerfile.builder")
    replace_in_file(dockerfile,
                    'version.AppName="maanyappd"',
                    f'version.AppName="{BINARY_NAME}"')
    replace_in_file(dockerfile,
                    '-o /neutron/build/maanyappd',
                    f'-o /neutron/build/{BINARY_NAME}')
    # If you also renamed the cmd folder, update the source path too
    replace_in_file(dockerfile,
                    '/neutron/cmd/maanyappd',
                    f'/neutron/cmd/{BINARY_NAME}')
    replace_in_file(dockerfile,
                    'COPY --from=builder /neutron/build/maanyappd /bin/maanyappd',
                    f'COPY --from=builder /neutron/build/{BINARY_NAME} /bin/{BINARY_NAME}')
    replace_in_file(dockerfile,
                    'ENTRYPOINT ["maanyappd"]',
                    f'ENTRYPOINT ["{BINARY_NAME}"]')

    # Scripts: defaults
    def shell_default(var: str, val: str) -> str:
        return '${' + var + ':-' + val + '}'

    replace_in_file(os.path.join(ROOT, 'network', 'init.sh'),
                    'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))
    replace_in_file(os.path.join(ROOT, 'network', 'init.sh'),
                    'CHAINID=${CHAINID:-test-1}', 'CHAINID=' + shell_default('CHAINID', CHAIN_ID))
    replace_in_file(os.path.join(ROOT, 'network', 'init.sh'),
                    'STAKEDENOM=${STAKEDENOM:-uapp}', 'STAKEDENOM=' + shell_default('STAKEDENOM', BASE_DENOM))

    replace_in_file(os.path.join(ROOT, 'network', 'init-neutrond.sh'),
                    'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))
    replace_in_file(os.path.join(ROOT, 'network', 'init-neutrond.sh'),
                    'CHAINID=${CHAINID:-test-1}', 'CHAINID=' + shell_default('CHAINID', CHAIN_ID))
    replace_in_file(os.path.join(ROOT, 'network', 'init-neutrond.sh'),
                    'STAKEDENOM=${STAKEDENOM:-uapp}', 'STAKEDENOM=' + shell_default('STAKEDENOM', BASE_DENOM))
    replace_in_file(os.path.join(ROOT, 'network', 'init-neutrond.sh'),
                    '"denom":"uapp"', f'"denom":"{BASE_DENOM}"')

    replace_in_file(os.path.join(ROOT, 'network', 'start.sh'),
                    'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))

    replace_in_file(os.path.join(ROOT, 'contrib', 'statesync.bash'),
                    'maanyappd init test', f'{BINARY_NAME} init test')
    replace_in_file(os.path.join(ROOT, 'contrib', 'statesync.bash'),
                    'minimum-gas-prices 0uapp', f'minimum-gas-prices 0{BASE_DENOM}')

    replace_in_file(os.path.join(ROOT, 'local_net', 'README.md'),
                    '--chain-id maanyapp-local-1', f'--chain-id {CHAIN_ID}')
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'single-val.sh'),
                    'maanyapp-local-1', CHAIN_ID)
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'single-val.sh'),
                    'maanyappd', BINARY_NAME)
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'utils.sh'),
                    '"denom": "uapp"', f'"denom": "{BASE_DENOM}"')
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'utils.sh'),
                    '"denom": "APP"', f'"denom": "{HUMAN_DENOM}"')
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'utils.sh'),
                    '"display": "APP"', f'"display": "{HUMAN_DENOM}"')
    replace_in_file(os.path.join(ROOT, 'local_net', 'sh', 'utils.sh'),
                    '"symbol": "APP"', f'"symbol": "{HUMAN_DENOM}"')
    # Optionally rename cmd folder if it exists
    old_cmd = os.path.join(ROOT, 'cmd', 'maanyappd')
    new_cmd = os.path.join(ROOT, 'cmd', BINARY_NAME)
    try:
        if os.path.exists(old_cmd) and old_cmd != new_cmd:
            os.rename(old_cmd, new_cmd)
            # Update references in Makefile to new cmd path
            replace_in_file(makefile, './cmd/maanyappd', f'./cmd/{BINARY_NAME}')
    except Exception:
        pass

    return 0


if __name__ == "__main__":
    sys.exit(main())
