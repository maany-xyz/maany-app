#!/usr/bin/env python3
"""
Cookiecutter post-gen for subdirectory template.

This hook copies a source codebase into the newly generated project directory
and then applies parameter substitutions. To locate the source, set env var

    SOURCE_DIR=/absolute/path/to/maany-app

before running cookiecutter with --directory template/cc.
"""

import os
import shutil
import sys


def is_text_file(path: str) -> bool:
    try:
        with open(path, "rb") as f:
            chunk = f.read(4096)
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


def copy_source_tree(src_root: str, dst_root: str) -> None:
    exclude_dirs = {".git", "build", "release", "vendor", "template"}
    for dirpath, dirnames, filenames in os.walk(src_root):
        # Prune excluded directories
        dirnames[:] = [d for d in dirnames if d not in exclude_dirs]
        rel = os.path.relpath(dirpath, src_root)
        out_dir = os.path.join(dst_root, rel) if rel != "." else dst_root
        os.makedirs(out_dir, exist_ok=True)
        for name in filenames:
            # Skip cookiecutter-specific files to avoid leaking template internals
            if name.endswith(".pyc"):
                continue
            if name == "cookiecutter.json":
                continue
            src = os.path.join(dirpath, name)
            dst = os.path.join(out_dir, name)
            shutil.copy2(src, dst)


def walk_and_replace(root: str, before: str, after: str) -> None:
    for dirpath, dirnames, filenames in os.walk(root):
        if ".git" in dirnames:
            dirnames.remove(".git")
        for name in filenames:
            path = os.path.join(dirpath, name)
            if not is_text_file(path):
                continue
            replace_in_file(path, before, after)


def main() -> int:
    dst_root = os.getcwd()
    src_root = os.environ.get("SOURCE_DIR", "")
    if not src_root or not os.path.isdir(src_root):
        print("ERROR: SOURCE_DIR env var must be set to the base repo path.", file=sys.stderr)
        return 2

    # 1) Copy source tree into this generated directory
    copy_source_tree(src_root, dst_root)

    # 2) Apply substitutions using the same logic as the top-level hook
    DEFAULT_MODULE = "github.com/maany-xyz/maany-app"
    TARGET_MODULE = "{{ cookiecutter.module_path }}"
    BINARY_NAME = "{{ cookiecutter.binary_name }}"
    BECH32_PREFIX = "{{ cookiecutter.bech32_main_prefix }}"
    HUMAN_DENOM = "{{ cookiecutter.display_denom }}"
    BASE_DENOM = "{{ cookiecutter.base_denom }}"
    DENOM_EXP = "{{ cookiecutter.denom_exponent }}"
    CHAIN_ID = "{{ cookiecutter.chain_id }}"

    # Replace module root and go.mod
    walk_and_replace(dst_root, DEFAULT_MODULE, TARGET_MODULE)
    replace_in_file(os.path.join(dst_root, "go.mod"), f"module {DEFAULT_MODULE}", f"module {TARGET_MODULE}")

    # Core params in code
    replace_in_file(os.path.join(dst_root, "app", "app.go"),
                    'Name = "maanyappd"', f'Name = "{BINARY_NAME}"')
    cfg_path = os.path.join(dst_root, "app", "config", "config.go")
    replace_in_file(cfg_path, 'HumanCoinUnit     = "APP"', f'HumanCoinUnit     = "{HUMAN_DENOM}"')
    replace_in_file(cfg_path, 'BaseCoinUnit      = "uapp"', f'BaseCoinUnit      = "{BASE_DENOM}"')
    replace_in_file(cfg_path, 'DefaultBondDenom  = BaseCoinUnit', 'DefaultBondDenom  = BaseCoinUnit')
    replace_in_file(cfg_path, 'DefaultExponent   = 6', f'DefaultExponent   = {DENOM_EXP}')
    replace_in_file(cfg_path, 'Bech32MainPrefix = "maanyapp"', f'Bech32MainPrefix = "{BECH32_PREFIX}"')
    replace_in_file(os.path.join(dst_root, "app", "params", "denom.go"),
                    'const DefaultDenom = "uapp"', f'const DefaultDenom = "{BASE_DENOM}"')
    replace_in_file(os.path.join(dst_root, "app", "genesis.go"),
                    'var FeeDenom = "uapp"', f'var FeeDenom = "{BASE_DENOM}"')

    # Makefile tweaks
    makefile = os.path.join(dst_root, "Makefile")
    replace_in_file(makefile, 'version.Name=maanyapp', 'version.Name={{ cookiecutter.project_name | replace(" ", "") | lower }}')
    replace_in_file(makefile, 'version.AppName=maanyappd', f'version.AppName={BINARY_NAME}')
    replace_in_file(makefile, 'Starting up maanyappd alone...', f'Starting up {BINARY_NAME} alone...')
    replace_in_file(makefile, 'export BINARY=maanyappd', f'export BINARY={BINARY_NAME}')
    replace_in_file(makefile, 'CHAINID=test-1', f'CHAINID={CHAIN_ID}')
    replace_in_file(makefile, 'STAKEDENOM=uapp', f'STAKEDENOM={BASE_DENOM}')
    replace_in_file(makefile, 'Killing maanyappd', f'Killing {BINARY_NAME}')
    replace_in_file(makefile, 'BINARY=maanyappd .', f'BINARY={BINARY_NAME} .')

    # Dockerfile
    dockerfile = os.path.join(dst_root, "Dockerfile.builder")
    replace_in_file(dockerfile, 'version.AppName="maanyappd"', f'version.AppName="{BINARY_NAME}"')
    replace_in_file(dockerfile, '-o /neutron/build/maanyappd', f'-o /neutron/build/{BINARY_NAME}')
    replace_in_file(dockerfile, '/neutron/cmd/maanyappd', f'/neutron/cmd/{BINARY_NAME}')
    replace_in_file(dockerfile, 'COPY --from=builder /neutron/build/maanyappd /bin/maanyappd', f'COPY --from=builder /neutron/build/{BINARY_NAME} /bin/{BINARY_NAME}')
    replace_in_file(dockerfile, 'ENTRYPOINT ["maanyappd"]', f'ENTRYPOINT ["{BINARY_NAME}"]')

    # Scripts: defaults (avoid Jinja braces in this hook source)
    def shell_default(var: str, val: str) -> str:
        return '${' + var + ':-' + val + '}'

    replace_in_file(os.path.join(dst_root, 'network', 'init.sh'), 'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))
    replace_in_file(os.path.join(dst_root, 'network', 'init.sh'), 'CHAINID=${CHAINID:-test-1}', 'CHAINID=' + shell_default('CHAINID', CHAIN_ID))
    replace_in_file(os.path.join(dst_root, 'network', 'init.sh'), 'STAKEDENOM=${STAKEDENOM:-uapp}', 'STAKEDENOM=' + shell_default('STAKEDENOM', BASE_DENOM))
    replace_in_file(os.path.join(dst_root, 'network', 'init-neutrond.sh'), 'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))
    replace_in_file(os.path.join(dst_root, 'network', 'init-neutrond.sh'), 'CHAINID=${CHAINID:-test-1}', 'CHAINID=' + shell_default('CHAINID', CHAIN_ID))
    replace_in_file(os.path.join(dst_root, 'network', 'init-neutrond.sh'), 'STAKEDENOM=${STAKEDENOM:-uapp}', 'STAKEDENOM=' + shell_default('STAKEDENOM', BASE_DENOM))
    replace_in_file(os.path.join(dst_root, 'network', 'init-neutrond.sh'), '"denom":"uapp"', f'"denom":"{BASE_DENOM}"')
    replace_in_file(os.path.join(dst_root, 'network', 'start.sh'), 'BINARY=${BINARY:-maanyappd}', 'BINARY=' + shell_default('BINARY', BINARY_NAME))
    replace_in_file(os.path.join(dst_root, 'contrib', 'statesync.bash'), 'maanyappd init test', f'{BINARY_NAME} init test')
    replace_in_file(os.path.join(dst_root, 'contrib', 'statesync.bash'), 'minimum-gas-prices 0uapp', f'minimum-gas-prices 0{BASE_DENOM}')
    replace_in_file(os.path.join(dst_root, 'local_net', 'README.md'), '--chain-id maanyapp-local-1', f'--chain-id {CHAIN_ID}')
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'single-val.sh'), 'maanyapp-local-1', CHAIN_ID)
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'single-val.sh'), 'maanyappd', BINARY_NAME)
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'utils.sh'), '"denom": "uapp"', f'"denom": "{BASE_DENOM}"')
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'utils.sh'), '"denom": "APP"', f'"denom": "{HUMAN_DENOM}"')
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'utils.sh'), '"display": "APP"', f'"display": "{HUMAN_DENOM}"')
    replace_in_file(os.path.join(dst_root, 'local_net', 'sh', 'utils.sh'), '"symbol": "APP"', f'"symbol": "{HUMAN_DENOM}"')

    # Optionally rename cmd folder
    old_cmd = os.path.join(dst_root, 'cmd', 'maanyappd')
    new_cmd = os.path.join(dst_root, 'cmd', BINARY_NAME)
    try:
        if os.path.exists(old_cmd) and old_cmd != new_cmd:
            os.rename(old_cmd, new_cmd)
            replace_in_file(makefile, './cmd/maanyappd', f'./cmd/{BINARY_NAME}')
    except Exception:
        pass

    return 0


if __name__ == "__main__":
    sys.exit(main())
