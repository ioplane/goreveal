from __future__ import annotations

import json
import re
from typing import Any


def normalize_goresym(raw: dict[str, Any]) -> dict[str, Any]:
    functions = sorted(
        {
            fn.get("FullName", "")
            for fn in raw.get("UserFunctions", [])
            if isinstance(fn, dict) and fn.get("FullName", "")
        }
    )

    return {
        "build_info": {
            "go_version": raw.get("BuildInfo", {}).get("GoVersion", ""),
            "path": raw.get("BuildInfo", {}).get("Path", ""),
        },
        "files": raw.get("Files", []),
        "functions": functions,
    }


def normalize_redress(
    packages_output: str,
    source_output: str,
    gomod_output: str = "",
) -> dict[str, Any]:
    packages: list[str] = []
    for raw_line in packages_output.splitlines():
        line = raw_line.strip()
        if not line or line in {"Packages:", "Name  Version", "----  -------"}:
            continue
        parts = line.split()
        if parts:
            packages.append(parts[0])

    source_files: list[str] = []
    functions: list[str] = []
    current_package = ""

    for raw_line in source_output.splitlines():
        line = raw_line.strip()
        if not line:
            continue
        if line.startswith("Package "):
            match = re.match(r"^Package\s+([^:]+):", line)
            if match:
                current_package = match.group(1).strip()
            continue
        if line.startswith("File: "):
            source_files.append(line.removeprefix("File: ").strip())
            continue
        if " Lines: " in line:
            name = line.split(" Lines: ", 1)[0].strip()
            if not name:
                continue
            if current_package and "." not in name and "/" not in name:
                name = f"{current_package}.{name}"
            functions.append(name)

    module_path = ""
    for raw_line in gomod_output.splitlines():
        line = raw_line.strip()
        if not line or line.startswith("Type ") or set(line) == {"-"}:
            continue
        parts = line.split()
        if len(parts) >= 2 and parts[0] == "main":
            module_path = parts[1]
            break

    return {
        "build_info": {
            "path": module_path,
        },
        "packages": sorted(set(packages)),
        "source_files": sorted(set(source_files)),
        "functions": sorted(set(functions)),
    }


def normalize_gore(raw: dict[str, Any]) -> dict[str, Any]:
    return {
        "build_info": {
            "go_version": raw.get("build_info", {}).get("go_version", ""),
            "path": raw.get("build_info", {}).get("path", ""),
        },
        "packages": sorted(set(raw.get("packages", []))),
        "source_files": sorted(set(raw.get("source_files", []))),
        "functions": sorted(set(raw.get("functions", []))),
        "types": sorted(set(raw.get("types", []))),
    }


def dump_json(payload: dict[str, Any]) -> str:
    return json.dumps(payload, indent=2)
