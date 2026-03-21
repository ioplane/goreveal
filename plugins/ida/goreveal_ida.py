"""Thin IDA adapter for GoREveal export payloads.

This module intentionally keeps all recovery logic outside of the plugin layer.
It only validates the export contract and turns payloads into applyable actions.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Any


IDA_CONTRACT_V1 = "goreveal.export.ida/v1"


def load_export(source: str | Path | dict[str, Any]) -> dict[str, Any]:
    if isinstance(source, dict):
        payload = source
    else:
        text = str(source)
        stripped = text.lstrip()
        if stripped.startswith("{") or stripped.startswith("["):
            payload = json.loads(text)
        else:
            payload = json.loads(Path(source).read_text(encoding="utf-8"))

    if payload.get("contract") != IDA_CONTRACT_V1:
        raise ValueError(f"unsupported contract: {payload.get('contract')!r}")

    return payload


def build_actions(payload: dict[str, Any]) -> list[dict[str, Any]]:
    payload = load_export(payload)
    actions: list[dict[str, Any]] = []

    for function in payload.get("functions", []):
        actions.append(
            {
                "kind": "function",
                "start": function["entry"],
                "end": function["end"],
                "name": function.get("refined_name") or function["name"],
                "raw_name": function["name"],
            }
        )

    for typ in payload.get("types", []):
        actions.append(
            {
                "kind": "type",
                "name": typ.get("refined_name") or typ["name"],
                "raw_name": typ["name"],
                "type_kind": typ["kind"],
            }
        )

    for string in payload.get("strings", []):
        actions.append(
            {
                "kind": "string",
                "value": string.get("refined_value") or string["value"],
                "raw_value": string["value"],
                "address": string.get("address"),
                "offset": string["offset"],
                "region": string["region"],
            }
        )

    for package in payload.get("packages", []):
        actions.append(
            {
                "kind": "package",
                "name": package["name"],
                "function_count": package.get("function_count", 0),
            }
        )

    source_tree = payload.get("source_tree") or {}
    for file_name in source_tree.get("files", []):
        actions.append(
            {
                "kind": "source_file",
                "name": file_name,
                "root": source_tree.get("root", ""),
            }
        )

    return actions


def apply_actions(api: Any, actions: list[dict[str, Any]]) -> None:
    for action in actions:
        kind = action["kind"]
        if kind == "function":
            api.ensure_function(action["start"], action["end"])
            api.rename_function(action["start"], action["name"])
        elif kind == "package":
            api.ensure_folder(action["name"])
        elif kind == "type":
            api.set_comment(0, f"type {action['name']} ({action['type_kind']})")
        elif kind == "string":
            api.set_comment(action.get("address") or action["offset"], action["value"])


def main(argv: list[str] | None = None) -> int:
    import argparse

    parser = argparse.ArgumentParser(description="Validate and stage GoREveal IDA export actions")
    parser.add_argument("input", help="path to goreveal export ida JSON")
    args = parser.parse_args(argv)

    payload = load_export(args.input)
    actions = build_actions(payload)
    print(json.dumps({"contract": IDA_CONTRACT_V1, "actions": actions}, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
