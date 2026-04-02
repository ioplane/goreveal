from __future__ import annotations

import json
from typing import Any


def _contains_package(packages: list[dict[str, Any]], want: str) -> bool:
    return any(pkg.get("name") == want for pkg in packages if isinstance(pkg, dict))


def _add_check(
    checks: list[dict[str, Any]],
    *,
    baseline: str,
    subject: str,
    status: str,
    expected: Any = None,
    actual: Any = None,
    note: str = "",
) -> None:
    entry: dict[str, Any] = {
        "baseline": baseline,
        "subject": subject,
        "status": status,
    }
    if expected is not None:
        entry["expected"] = expected
    if actual is not None:
        entry["actual"] = actual
    if note:
        entry["note"] = note
    checks.append(entry)


def build_fixture_report(
    *,
    fixture: str,
    analysis: dict[str, Any],
    goresym: dict[str, Any],
    redress: dict[str, Any],
    gore: dict[str, Any],
) -> dict[str, Any]:
    checks: list[dict[str, Any]] = []

    analysis_build = analysis.get("build_info", {})
    analysis_path = analysis_build.get("path", "")
    analysis_go_version = analysis_build.get("go_version", "")
    source_tree = analysis.get("source_tree") or {}
    source_root = source_tree.get("root", "")
    source_files = source_tree.get("files") or []
    first_file = source_files[0] if source_files else ""
    goresym_expected_file = f"{source_root}/{first_file}" if source_root and first_file else ""

    evidence_checks = [
        (
            "goresym",
            "build_info.path",
            analysis_path,
            goresym.get("build_info", {}).get("path", ""),
        ),
        (
            "goresym",
            "build_info.go_version",
            analysis_go_version,
            goresym.get("build_info", {}).get("go_version", ""),
        ),
        (
            "redress",
            "build_info.path",
            analysis_path,
            redress.get("build_info", {}).get("path", ""),
        ),
        (
            "gore",
            "build_info.path",
            analysis_path,
            gore.get("build_info", {}).get("path", ""),
        ),
        (
            "gore",
            "build_info.go_version",
            analysis_go_version,
            gore.get("build_info", {}).get("go_version", ""),
        ),
        (
            "goresym",
            "files.module_local",
            goresym_expected_file,
            goresym_expected_file if goresym_expected_file in goresym.get("files", []) else "",
        ),
        (
            "redress",
            "packages.main",
            True,
            "main" in redress.get("packages", []),
        ),
        (
            "gore",
            "packages.main",
            True,
            "main" in gore.get("packages", []),
        ),
        (
            "goreveal",
            "packages.main",
            True,
            _contains_package(analysis.get("packages", []), "main"),
        ),
        (
            "redress",
            "source_files.main.go",
            first_file,
            first_file if first_file in redress.get("source_files", []) else "",
        ),
        (
            "gore",
            "source_files.main.go",
            first_file,
            first_file if first_file in gore.get("source_files", []) else "",
        ),
        (
            "goresym",
            "functions.main.main",
            True,
            "main.main" in goresym.get("functions", []),
        ),
        (
            "goresym",
            "functions.main.helperAdd",
            True,
            "main.helperAdd" in goresym.get("functions", []),
        ),
        (
            "goresym",
            "functions.main.helperBanner",
            True,
            "main.helperBanner" in goresym.get("functions", []),
        ),
        (
            "redress",
            "functions.main.main",
            True,
            "main.main" in redress.get("functions", []),
        ),
        (
            "redress",
            "functions.main.helperBanner",
            True,
            "main.helperBanner" in redress.get("functions", []),
        ),
        (
            "gore",
            "functions.main.main",
            True,
            "main.main" in gore.get("functions", []),
        ),
        (
            "gore",
            "functions.main.helperBanner",
            True,
            "main.helperBanner" in gore.get("functions", []),
        ),
    ]

    for baseline, subject, expected, actual in evidence_checks:
        status = "matched" if expected == actual else "diverged"
        _add_check(
            checks,
            baseline=baseline,
            subject=subject,
            status=status,
            expected=expected,
            actual=actual,
        )

    _add_check(
        checks,
        baseline="redress",
        subject="build_info.go_version",
        status="skipped",
        note="Current redress evidence only covers module path via gomod, not go version parity.",
    )
    _add_check(
        checks,
        baseline="gore",
        subject="types.user_defined",
        status="skipped",
        note=(
            "Current gore type surface is runtime-heavy and not yet a trustworthy "
            "user-type parity surface."
        ),
    )
    _add_check(
        checks,
        baseline="goresym",
        subject="functions.complete_parity",
        status="skipped",
        note=(
            "Current GoReSym evidence intentionally covers a narrow "
            "user-function subset only."
        ),
    )

    summary = {
        "matched": sum(1 for entry in checks if entry["status"] == "matched"),
        "diverged": sum(1 for entry in checks if entry["status"] == "diverged"),
        "skipped": sum(1 for entry in checks if entry["status"] == "skipped"),
    }

    return {
        "fixture": fixture,
        "summary": summary,
        "checks": checks,
    }


def dump_report(report: dict[str, Any]) -> str:
    return json.dumps(report, indent=2)
