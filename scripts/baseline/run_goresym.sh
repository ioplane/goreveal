#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: run_goresym.sh <binary>" >&2
  exit 2
fi

BASELINES_ROOT="${GOREVEAL_BASELINES_ROOT:-/repos}"
BINARY_PATH="$(realpath "$1")"
PROJECT_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
export PYTHONPATH="${PROJECT_ROOT}${PYTHONPATH:+:${PYTHONPATH}}"

TMP_JSON="$(mktemp)"
TMP_ERR="$(mktemp)"
trap 'rm -f "$TMP_JSON" "$TMP_ERR"' EXIT

if ! (
  cd "${BASELINES_ROOT}/GoReSym"
  /usr/local/go/bin/go run . -p "${BINARY_PATH}"
) >"${TMP_JSON}" 2>"${TMP_ERR}"; then
  cat "${TMP_ERR}" >&2
  exit 1
fi

python3 - "${TMP_JSON}" <<'PY'
import json
import sys

from scripts.baseline.normalize import dump_json, normalize_goresym

with open(sys.argv[1], "r", encoding="utf-8") as fh:
    raw = json.load(fh)

print(dump_json(normalize_goresym(raw)))
PY
