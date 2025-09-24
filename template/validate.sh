#!/usr/bin/env bash
set -euo pipefail

# Validate Cookiecutter context against CUE schema if cue is installed.
if ! command -v cue >/dev/null 2>&1; then
  echo "cue not installed; skipping validation" >&2
  exit 0
fi

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

# Rendered cookiecutter.json acts as context here. In cookiecutter hooks,
# you can pass the context file; for source repo validation, this is a hint.
cue vet "${ROOT_DIR}/cookiecutter.json" "${ROOT_DIR}/template/SCHEMA.cue"
echo "CUE validation passed"

