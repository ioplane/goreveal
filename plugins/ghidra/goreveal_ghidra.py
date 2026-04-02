"""Thin Ghidra adapter for GoREveal export payloads.

The adapter consumes stable export data and stages import actions.
It does not perform recovery or deobfuscation on its own.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any

GHIDRA_CONTRACT_V1 = "goreveal.export.ghidra/v1"


def load_export(source: str | Path | dict[str, Any]) -> dict[str, Any]:
    if isinstance(source, dict):
        payload = source
    else:
        text = str(source)
        stripped = text.lstrip()
        if stripped.startswith(("{", "[")):
            payload = json.loads(text)
        else:
            payload = json.loads(Path(source).read_text(encoding="utf-8"))

    if payload.get("contract") != GHIDRA_CONTRACT_V1:
        raise ValueError(f"unsupported contract: {payload.get('contract')!r}")

    return payload


def build_actions(payload: dict[str, Any]) -> list[dict[str, Any]]:
    payload = load_export(payload)
    actions: list[dict[str, Any]] = []

    program = payload.get("program") or {}
    actions.append(
        {
            "kind": "program",
            "path": program.get("path", ""),
            "module_path": program.get("module_path", ""),
            "go_version": program.get("go_version", ""),
        }
    )

    actions.extend(
        {
            "kind": "symbol",
            "address": symbol["address"],
            "end": symbol["end"],
            "name": symbol.get("refined_name") or symbol["name"],
            "raw_name": symbol["name"],
        }
        for symbol in payload.get("symbols", [])
    )

    actions.extend(
        {
            "kind": "package",
            "name": package["name"],
            "function_count": package.get("function_count", 0),
        }
        for package in payload.get("packages", [])
    )

    actions.extend(
        {
            "kind": "type",
            "name": typ.get("refined_name") or typ["name"],
            "raw_name": typ["name"],
            "type_kind": typ["kind"],
        }
        for typ in payload.get("types", [])
    )

    actions.extend(
        {
            "kind": "string",
            "value": string.get("refined_value") or string["value"],
            "raw_value": string["value"],
            "address": string.get("address"),
            "offset": string["offset"],
            "region": string["region"],
        }
        for string in payload.get("strings", [])
    )

    source_tree = payload.get("source_tree") or {}
    actions.extend(
        {
            "kind": "source_file",
            "name": file_name,
            "root": source_tree.get("root", ""),
        }
        for file_name in source_tree.get("files", [])
    )

    return actions


def apply_actions(api: Any, actions: list[dict[str, Any]]) -> None:
    for action in actions:
        kind = action["kind"]
        if kind == "program":
            api.set_program_context(action["path"], action["module_path"], action["go_version"])
        elif kind == "symbol":
            api.ensure_function(action["address"], action["end"])
            api.rename_symbol(action["address"], action["name"])
        elif kind == "package":
            api.ensure_namespace(action["name"])
        elif kind == "type":
            api.set_comment(0, f"type {action['name']} ({action['type_kind']})")
        elif kind == "string":
            api.set_comment(action.get("address") or action["offset"], action["value"])


def main(argv: list[str] | None = None) -> int:
    import argparse

    parser = argparse.ArgumentParser(
        description="Validate and stage GoREveal Ghidra export actions"
    )
    parser.add_argument("input", help="path to goreveal export ghidra JSON")
    args = parser.parse_args(argv)

    payload = load_export(args.input)
    actions = build_actions(payload)
    print(json.dumps({"contract": GHIDRA_CONTRACT_V1, "actions": actions}, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
