#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: run_redress.sh <binary>" >&2
  exit 2
fi

BASELINES_ROOT="${GOREVEAL_BASELINES_ROOT:-/repos}"
BINARY_PATH="$(realpath "$1")"
PROJECT_ROOT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/../.." && pwd)"
export PYTHONPATH="${PROJECT_ROOT}${PYTHONPATH:+:${PYTHONPATH}}"

TMP_TXT="$(mktemp)"
TMP_SRC="$(mktemp)"
TMP_GOMOD="$(mktemp)"
TMP_ERR="$(mktemp)"
trap 'rm -f "$TMP_TXT" "$TMP_SRC" "$TMP_GOMOD" "$TMP_ERR"' EXIT

if ! (
  cd "${BASELINES_ROOT}/redress"
  /usr/local/go/bin/go run . packages "${BINARY_PATH}"
) >"${TMP_TXT}" 2>"${TMP_ERR}"; then
  cat "${TMP_ERR}" >&2
  exit 1
fi

if ! (
  cd "${BASELINES_ROOT}/redress"
  /usr/local/go/bin/go run . source "${BINARY_PATH}"
) >"${TMP_SRC}" 2>"${TMP_ERR}"; then
  cat "${TMP_ERR}" >&2
  exit 1
fi

if ! (
  cd "${BASELINES_ROOT}/redress"
  /usr/local/go/bin/go run . gomod "${BINARY_PATH}"
) >"${TMP_GOMOD}" 2>"${TMP_ERR}"; then
  cat "${TMP_ERR}" >&2
  exit 1
fi

python3 - "${TMP_TXT}" "${TMP_SRC}" "${TMP_GOMOD}" <<'PY'
import sys

from scripts.baseline.normalize import dump_json, normalize_redress

with open(sys.argv[1], "r", encoding="utf-8") as fh:
    packages_output = fh.read()
with open(sys.argv[2], "r", encoding="utf-8") as fh:
    source_output = fh.read()
with open(sys.argv[3], "r", encoding="utf-8") as fh:
    gomod_output = fh.read()

print(dump_json(normalize_redress(packages_output, source_output, gomod_output)))
PY
