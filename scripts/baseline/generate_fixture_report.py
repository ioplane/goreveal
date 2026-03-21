from __future__ import annotations

import json
import os
import subprocess
import sys
from pathlib import Path

from scripts.baseline.report import build_fixture_report, dump_report


def _run_json(cmd: list[str], *, cwd: Path, env: dict[str, str]) -> dict:
    completed = subprocess.run(
        cmd,
        cwd=cwd,
        env=env,
        text=True,
        capture_output=True,
        check=True,
    )
    return json.loads(completed.stdout)


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: python -m scripts.baseline.generate_fixture_report <fixture>", file=sys.stderr)
        return 2

    fixture = sys.argv[1]
    project_root = Path(__file__).resolve().parents[2]
    env = os.environ.copy()
    goreveal_bin = env.get("GOREVEAL_BIN")
    analysis_cmd = [goreveal_bin, "analyze", fixture] if goreveal_bin else [
        "/usr/local/go/bin/go",
        "run",
        "./cmd/goreveal",
        "analyze",
        fixture,
    ]

    analysis = _run_json(
        analysis_cmd,
        cwd=project_root,
        env=env,
    )
    goresym = _run_json(
        [str(project_root / "scripts" / "baseline" / "run_goresym.sh"), fixture],
        cwd=project_root,
        env=env,
    )
    redress = _run_json(
        [str(project_root / "scripts" / "baseline" / "run_redress.sh"), fixture],
        cwd=project_root,
        env=env,
    )
    gore = _run_json(
        [str(project_root / "scripts" / "baseline" / "run_gore.sh"), fixture],
        cwd=project_root,
        env=env,
    )

    fixture_path = Path(fixture)
    fixture_name = fixture_path.parent.name if fixture_path.stem == "fixture" else fixture_path.stem

    report = build_fixture_report(
        fixture=fixture_name or fixture,
        analysis=analysis,
        goresym=goresym,
        redress=redress,
        gore=gore,
    )
    print(dump_report(report))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
