from __future__ import annotations

import argparse
import json
import os
import shutil
import subprocess
from dataclasses import dataclass
from pathlib import Path
from typing import Any

REPO_ROOT = Path(__file__).resolve().parents[2]
DEFAULT_SOURCE = REPO_ROOT / "corpus" / "protected" / "enterprise-sample" / "src"
DEFAULT_OUTPUT_DIR = REPO_ROOT / ".tmp" / "protected-matrix"
DEFAULT_GOREVEAL_BIN = REPO_ROOT / ".tmp" / "goreveal"
DEFAULT_GARBLE_SOURCE = Path(os.environ.get("GARBLE_SOURCE", "/repos/garble"))
RELEVANT_FUNCTIONS = (
    "main.readLicenseToken",
    "main.auditFeatureGate",
    "main.runEnterpriseReport",
)


@dataclass(frozen=True)
class BuildProfile:
    name: str
    supported_goos: tuple[str, ...]
    buildmode: str = ""
    trimpath: bool = False
    ldflags: str = ""
    garble_flags: tuple[str, ...] = ()
    gogarble: str = ""


@dataclass(frozen=True)
class PlatformTarget:
    goos: str
    goarch: str
    binary_name: str

    @property
    def name(self) -> str:
        return f"{self.goos}-{self.goarch}"


PROFILES: tuple[BuildProfile, ...] = (
    BuildProfile(name="plain", supported_goos=("linux", "windows", "darwin")),
    BuildProfile(
        name="stripped",
        supported_goos=("linux", "windows", "darwin"),
        ldflags="-s -w -buildid=",
    ),
    BuildProfile(name="trimpath", supported_goos=("linux", "windows", "darwin"), trimpath=True),
    BuildProfile(
        name="stripped-trimpath",
        supported_goos=("linux", "windows", "darwin"),
        trimpath=True,
        ldflags="-s -w -buildid=",
    ),
    BuildProfile(name="pie", supported_goos=("linux", "darwin"), buildmode="pie"),
    BuildProfile(
        name="garble",
        supported_goos=("linux",),
        garble_flags=(),
        gogarble="*",
    ),
    BuildProfile(
        name="garble-literals-tiny",
        supported_goos=("linux",),
        garble_flags=("-literals", "-tiny"),
        gogarble="*",
    ),
)

PLATFORMS: tuple[PlatformTarget, ...] = (
    PlatformTarget(goos="linux", goarch="amd64", binary_name="enterprise-sample"),
    PlatformTarget(goos="linux", goarch="arm64", binary_name="enterprise-sample"),
    PlatformTarget(goos="windows", goarch="amd64", binary_name="enterprise-sample.exe"),
    PlatformTarget(goos="windows", goarch="arm64", binary_name="enterprise-sample.exe"),
    PlatformTarget(goos="darwin", goarch="amd64", binary_name="enterprise-sample"),
    PlatformTarget(goos="darwin", goarch="arm64", binary_name="enterprise-sample"),
)


def build_go_binary(
    go_bin: str,
    garble_bin: str,
    source_dir: Path,
    output_path: Path,
    profile: BuildProfile,
    platform: PlatformTarget,
) -> None:
    build_args = ["build", "-o", str(output_path)]
    if profile.trimpath:
        build_args.append("-trimpath")
    if profile.buildmode:
        build_args.extend(["-buildmode", profile.buildmode])
    if profile.ldflags:
        build_args.extend(["-ldflags", profile.ldflags])
    build_args.append(".")

    env = os.environ.copy()
    env["CGO_ENABLED"] = "0"
    env["GOOS"] = platform.goos
    env["GOARCH"] = platform.goarch
    env["GOWORK"] = "off"
    toolchain_path = f"{Path(go_bin).parent}:{Path(garble_bin).parent}"
    if existing_path := env.get("PATH"):
        env["PATH"] = f"{toolchain_path}:{existing_path}"
    else:
        env["PATH"] = toolchain_path
    if profile.gogarble:
        env["GOGARBLE"] = profile.gogarble

    tool_args = [go_bin, *build_args]
    if profile.gogarble:
        tool_args = [garble_bin, *profile.garble_flags, *build_args]

    subprocess.run(
        tool_args,
        cwd=source_dir,
        env=env,
        text=True,
        capture_output=True,
        check=True,
    )


def ensure_garble_binary(
    *,
    go_bin: str,
    garble_bin: str,
    garble_source: Path | None,
    tools_dir: Path,
) -> str:
    if garble_source is None or not garble_source.exists():
        return garble_bin

    tools_dir.mkdir(parents=True, exist_ok=True)
    built_garble = tools_dir / "garble"
    env = os.environ.copy()
    env["GOWORK"] = "off"
    subprocess.run(
        [go_bin, "build", "-o", str(built_garble), "."],
        cwd=garble_source,
        env=env,
        text=True,
        capture_output=True,
        check=True,
    )
    return str(built_garble)


def _run_json(cmd: list[str], *, cwd: Path, env: dict[str, str]) -> dict[str, Any]:
    proc = subprocess.run(
        cmd,
        cwd=cwd,
        env=env,
        check=True,
        capture_output=True,
        text=True,
    )
    return json.loads(proc.stdout)


def _tool_function_presence(functions: list[str]) -> list[str]:
    return [name for name in RELEVANT_FUNCTIONS if name in functions]


def _count_source_files(analysis: dict[str, Any]) -> int:
    source_tree = analysis.get("source_tree") or {}
    return len(source_tree.get("files") or [])


def _display_path(path: Path) -> str:
    try:
        return str(path.relative_to(REPO_ROOT))
    except ValueError:
        return str(path)


def analyze_binary(goreveal_bin: Path, binary_path: Path) -> dict[str, Any]:
    env = os.environ.copy()
    analysis = _run_json(
        [str(goreveal_bin), "analyze", str(binary_path)],
        cwd=REPO_ROOT,
        env=env,
    )
    return {"analyze": analysis}


def _baseline_summary(wrapper: str, binary_path: Path) -> dict[str, Any]:
    payload = _run_json(
        ["bash", str(REPO_ROOT / "scripts" / "baseline" / wrapper), str(binary_path)],
        cwd=REPO_ROOT,
        env=os.environ.copy(),
    )
    return {
        "build_info_path": (payload.get("build_info") or {}).get("path", ""),
        "functions": len(payload.get("functions") or []),
        "packages": len(payload.get("packages") or []),
        "files": len(payload.get("files") or payload.get("source_files") or []),
        "relevant_functions": _tool_function_presence(payload.get("functions") or []),
    }


def _safe_baseline_summary(wrapper: str, binary_path: Path) -> dict[str, Any]:
    try:
        return _baseline_summary(wrapper, binary_path)
    except subprocess.CalledProcessError as exc:
        return {
            "error": exc.stderr.strip() or exc.stdout.strip() or "baseline wrapper failed",
        }


def profile_record(
    *,
    profile: BuildProfile,
    platform: PlatformTarget,
    binary_path: Path,
    analysis_bundle: dict[str, Any],
) -> dict[str, Any]:
    analysis = analysis_bundle["analyze"]
    build_info = analysis.get("build_info") or {}
    runtime = analysis.get("runtime") or {}
    source_tree = analysis.get("source_tree") or {}
    peeling = analysis.get("peeling") or {}
    functions = analysis.get("functions") or []
    peel_functions = peeling.get("functions") or []
    peel_packages = peeling.get("packages") or []

    goreveal = {
        "build_info_path": build_info.get("path", ""),
        "runtime_trust_summary": runtime.get("trust_summary", "absent"),
        "elf_text_section_addr": runtime.get("elf_text_section_addr", 0),
        "elf_text_section_end_inclusive": runtime.get("elf_text_section_end_inclusive", 0),
        "elf_pclntab_header_magic": runtime.get("elf_pclntab_header_magic", ""),
        "elf_pclntab_header_magic_kind": runtime.get("elf_pclntab_header_magic_kind", ""),
        "elf_pclntab_function_count_hint": runtime.get("elf_pclntab_function_count_hint", 0),
        "elf_pclntab_file_count_hint": runtime.get("elf_pclntab_file_count_hint", 0),
        "elf_functab_first_pc_offset_hint": runtime.get("elf_functab_first_pc_offset_hint", 0),
        "elf_functab_last_pc_offset_hint": runtime.get("elf_functab_last_pc_offset_hint", 0),
        "elf_functab_pc_offsets_monotonic": bool(
            runtime.get("elf_functab_pc_offsets_monotonic", False)
        ),
        "elf_functab_first_pc_addr_hint": runtime.get("elf_functab_first_pc_addr_hint", 0),
        "elf_functab_last_pc_addr_hint": runtime.get("elf_functab_last_pc_addr_hint", 0),
        "elf_functab_pc_addr_hints_within_text": bool(
            runtime.get("elf_functab_pc_addr_hints_within_text", False)
        ),
        "elf_functab_pc_addr_sample_count": len(runtime.get("elf_functab_pc_addr_sample") or []),
        "elf_functab_pc_addr_sample_first": (
            (runtime.get("elf_functab_pc_addr_sample") or [0])[0]
            if (runtime.get("elf_functab_pc_addr_sample") or [])
            else 0
        ),
        "elf_functab_pc_addr_sample_all_within_text": bool(
            runtime.get("elf_functab_pc_addr_sample_all_within_text", False)
        ),
        "elf_function_foothold": runtime.get("elf_function_foothold", ""),
        "elf_function_foothold_count_hint": runtime.get("elf_function_foothold_count_hint", 0),
        "elf_function_foothold_text_source": runtime.get("elf_function_foothold_text_source", ""),
        "elf_function_foothold_start_addr": runtime.get("elf_function_foothold_start_addr", 0),
        "elf_function_foothold_end_addr": runtime.get("elf_function_foothold_end_addr", 0),
        "elf_function_recovery_blocker": runtime.get("elf_function_recovery_blocker", ""),
        "moduledata_pcheader_matches_gopclntab": bool(
            runtime.get("moduledata_pcheader_matches_gopclntab", False)
        ),
        "moduledata_funcnametab_within_gopclntab": bool(
            runtime.get("moduledata_funcnametab_within_gopclntab", False)
        ),
        "moduledata_pclntable_within_gopclntab": bool(
            runtime.get("moduledata_pclntable_within_gopclntab", False)
        ),
        "functions": len(functions),
        "packages": len(analysis.get("packages") or []),
        "files": _count_source_files(analysis),
        "pathless_file_evidence": bool(source_tree.get("pathless_file_evidence", False)),
        "peeling_functions": len(peel_functions),
        "user_functions": sum(
            1 for fn in peel_functions if fn.get("classification") == "user"
        ),
        "user_packages": sum(
            1 for pkg in peel_packages if pkg.get("primary_classification") == "user"
        ),
        "relevant_functions": _tool_function_presence([fn.get("name", "") for fn in functions]),
    }

    return {
        "profile": profile.name,
        "target": platform.name,
        "binary": _display_path(binary_path),
        "buildmode": profile.buildmode or "default",
        "trimpath": profile.trimpath,
        "ldflags": profile.ldflags,
        "format": analysis["input"]["format"],
        "goreveal": goreveal,
        "goresym": _safe_baseline_summary("run_goresym.sh", binary_path),
        "redress": _safe_baseline_summary("run_redress.sh", binary_path),
        "gore": _safe_baseline_summary("run_gore.sh", binary_path),
    }


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Build and compare the first protected-binary profile matrix."
    )
    parser.add_argument("--source", type=Path, default=DEFAULT_SOURCE)
    parser.add_argument("--output-dir", type=Path, default=DEFAULT_OUTPUT_DIR)
    parser.add_argument(
        "--goreveal-bin",
        type=Path,
        default=Path(os.environ.get("GOREVEAL_BIN", str(DEFAULT_GOREVEAL_BIN))),
    )
    parser.add_argument("--go-bin", default="/usr/local/go/bin/go")
    parser.add_argument("--garble-bin", default="/go/bin/garble")
    parser.add_argument(
        "--garble-source",
        type=Path,
        default=Path(os.environ.get("GARBLE_SOURCE", str(DEFAULT_GARBLE_SOURCE))),
    )
    parser.add_argument("--json-out", type=Path)
    return parser.parse_args()


def generate_matrix(
    *,
    source_dir: Path,
    output_dir: Path,
    goreveal_bin: Path,
    go_bin: str,
    garble_bin: str,
    garble_source: Path | None,
) -> dict[str, Any]:
    if output_dir.exists():
        shutil.rmtree(output_dir)
    output_dir.mkdir(parents=True, exist_ok=True)

    if not source_dir.exists():
        raise SystemExit(f"source dir does not exist: {source_dir}")
    if not goreveal_bin.exists():
        raise SystemExit(f"goreveal binary does not exist: {goreveal_bin}")

    staged_source = output_dir / "source"
    shutil.copytree(source_dir, staged_source)
    resolved_garble_bin = ensure_garble_binary(
        go_bin=go_bin,
        garble_bin=garble_bin,
        garble_source=garble_source,
        tools_dir=output_dir / ".tools",
    )

    records: list[dict[str, Any]] = []
    skipped: list[dict[str, str]] = []
    for platform in PLATFORMS:
        for profile in PROFILES:
            if platform.goos not in profile.supported_goos:
                skipped.append(
                    {
                        "target": platform.name,
                        "profile": profile.name,
                        "reason": "profile not enabled for platform",
                    }
                )
                continue

            binary_path = (
                output_dir / platform.goos / platform.goarch / profile.name / platform.binary_name
            )
            binary_path.parent.mkdir(parents=True, exist_ok=True)
            try:
                build_go_binary(
                    go_bin,
                    resolved_garble_bin,
                    staged_source,
                    binary_path,
                    profile,
                    platform,
                )
            except subprocess.CalledProcessError as exc:
                skipped.append(
                    {
                        "target": platform.name,
                        "profile": profile.name,
                        "reason": exc.stderr.strip() or "build failed",
                    }
                )
                continue

            records.append(
                profile_record(
                    profile=profile,
                    platform=platform,
                    binary_path=binary_path,
                    analysis_bundle=analyze_binary(goreveal_bin, binary_path),
                )
            )

    return {
        "source": _display_path(source_dir),
        "goreveal_bin": _display_path(goreveal_bin),
        "garble_bin": _display_path(Path(resolved_garble_bin)),
        "garble_source": (
            _display_path(garble_source)
            if garble_source is not None and garble_source.exists()
            else ""
        ),
        "relevant_functions": list(RELEVANT_FUNCTIONS),
        "profiles": records,
        "skipped": skipped,
    }


def main() -> int:
    args = parse_args()
    payload = generate_matrix(
        source_dir=args.source.resolve(),
        output_dir=args.output_dir.resolve(),
        goreveal_bin=args.goreveal_bin.resolve(),
        go_bin=args.go_bin,
        garble_bin=args.garble_bin,
        garble_source=args.garble_source.resolve(),
    )
    if args.json_out is not None:
        args.json_out.parent.mkdir(parents=True, exist_ok=True)
        args.json_out.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")

    print(json.dumps(payload, indent=2))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
