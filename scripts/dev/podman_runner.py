from __future__ import annotations

import argparse
import os
import sys
from collections.abc import Mapping
from contextlib import suppress
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, NotRequired, TypedDict

DEFAULT_DEV_IMAGE = "localhost/goreveal:dev"
DEFAULT_DEV_WORKDIR = "/workspace"
# Host directory holding the reference-tool checkouts used by the differential
# suite. Override with GOREVEAL_BASELINES_HOST_ROOT; see CONTRIBUTING.md.
DEFAULT_BASELINES_ROOT = os.environ.get(
    "GOREVEAL_BASELINES_HOST_ROOT",
    str(Path.home() / "goreveal-baselines"),
)
GOFMT_ALL_CMD = '/usr/local/go/bin/gofmt -w $(find . -type f -name "*.go" -not -path "./.git/*")'


@dataclass(frozen=True)
class Step:
    cmd: list[str]
    with_baselines: bool = False
    env: dict[str, str] = field(default_factory=dict)
    image: str = DEFAULT_DEV_IMAGE


class MountSpec(TypedDict):
    bind: str
    mode: str


class RunKwargs(TypedDict):
    image: str
    command: list[str]
    working_dir: str
    volumes: dict[str, MountSpec]
    stderr: bool
    stdout: bool
    remove: bool
    environment: NotRequired[dict[str, str]]


def default_base_url(*, env: Mapping[str, str] | None = None, uid: int | None = None) -> str:
    current_env = dict(os.environ if env is None else env)
    explicit = current_env.get("PODMAN_BASE_URL")
    if explicit:
        return explicit

    container_host = current_env.get("CONTAINER_HOST")
    if container_host:
        return container_host

    docker_host = current_env.get("DOCKER_HOST")
    if docker_host:
        return docker_host

    runtime_dir = current_env.get("XDG_RUNTIME_DIR")
    if runtime_dir:
        return f"unix://{runtime_dir}/podman/podman.sock"

    resolved_uid = os.getuid() if uid is None else uid
    return f"unix:///run/user/{resolved_uid}/podman/podman.sock"


def create_client(base_url: str | None = None) -> Any:
    try:
        from podman import PodmanClient
    except ImportError as exc:  # pragma: no cover - exercised only when dependency is missing
        raise SystemExit(
            "podman Python package is not installed. Install project Python dependencies first."
        ) from exc

    return PodmanClient(base_url=base_url or default_base_url())


def run_kwargs_for_step(step: Step, *, cwd: str) -> RunKwargs:
    workspace = str(Path(cwd).resolve())
    volumes: dict[str, MountSpec] = {
        workspace: {"bind": DEFAULT_DEV_WORKDIR, "mode": "Z"},
    }
    environment = dict(step.env)

    if step.with_baselines:
        volumes[DEFAULT_BASELINES_ROOT] = {"bind": "/repos", "mode": "Z"}
        environment["GOREVEAL_BASELINES_ROOT"] = "/repos"

    kwargs: RunKwargs = {
        "image": step.image,
        "command": step.cmd,
        "working_dir": DEFAULT_DEV_WORKDIR,
        "volumes": volumes,
        "stderr": True,
        "stdout": True,
        "remove": False,
    }
    if environment:
        kwargs["environment"] = environment
    return kwargs


def build_dev_image(*, client: Any, repo_root: str, image: str = DEFAULT_DEV_IMAGE) -> None:
    image_tag = image
    context_path = str(Path(repo_root).resolve())
    dockerfile = str(Path("deployments") / "docker" / "Containerfile.dev")
    client.images.build(path=context_path, dockerfile=dockerfile, tag=image_tag)


def execute_step(*, client: Any, step: Step, cwd: str) -> int:
    kwargs = run_kwargs_for_step(step, cwd=cwd)
    container = client.containers.create(**kwargs)
    try:
        container.start()
        for chunk in container.logs(stream=True, stdout=True, stderr=True):
            if isinstance(chunk, bytes):
                sys.stdout.buffer.write(chunk)
            else:
                sys.stdout.write(str(chunk))
        result = container.wait()
        return int(result.get("StatusCode", 1)) if isinstance(result, dict) else int(result)
    finally:
        with suppress(Exception):
            container.remove(force=True)


def task_steps(name: str) -> list[Step]:
    workspace_bin = f"{DEFAULT_DEV_WORKDIR}/.tmp/goreveal"
    protected_matrix_bin = f"{DEFAULT_DEV_WORKDIR}/.tmp/protected/goreveal"
    tasks: dict[str, list[Step]] = {
        "fmt": [
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    GOFMT_ALL_CMD,
                ]
            ),
        ],
        "build-dev-bin": [
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    "mkdir -p .tmp && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal",
                ]
            ),
        ],
        "test": [
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    "mkdir -p .tmp && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal",
                ]
            ),
            Step(cmd=["/usr/local/go/bin/go", "test", "./..."], with_baselines=True),
            Step(cmd=["python3", "-m", "unittest", "scripts.baseline.test_normalize"]),
            Step(cmd=["python3", "-m", "unittest", "scripts.baseline.test_report"]),
            Step(
                cmd=["python3", "-m", "unittest", "plugins.ida.test_goreveal_ida"],
                env={"GOREVEAL_BIN": workspace_bin},
            ),
            Step(
                cmd=["python3", "-m", "unittest", "plugins.ghidra.test_goreveal_ghidra"],
                env={"GOREVEAL_BIN": workspace_bin},
            ),
        ],
        "test-plugins": [
            Step(
                cmd=["python3", "-m", "unittest", "plugins.ida.test_goreveal_ida"],
                env={"GOREVEAL_BIN": workspace_bin},
            ),
            Step(
                cmd=["python3", "-m", "unittest", "plugins.ghidra.test_goreveal_ghidra"],
                env={"GOREVEAL_BIN": workspace_bin},
            ),
        ],
        "test-differential": [
            Step(cmd=["/usr/local/go/bin/go", "test", "./tests/differential"], with_baselines=True),
        ],
        "test-differential-report": [
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    "mkdir -p .tmp && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal",
                ]
            ),
            Step(
                cmd=[
                    "python3",
                    "-m",
                    "scripts.baseline.generate_fixture_report",
                    "corpus/fixtures/go-elf-buildinfo-linux-amd64/fixture.bin",
                ],
                with_baselines=True,
                env={"GOREVEAL_BIN": workspace_bin},
            ),
        ],
        "protected-matrix": [
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    (
                        "cd /workspace && "
                        "mkdir -p .tmp/protected /tmp/protected-matrix && "
                        f"/usr/local/go/bin/go build -o {protected_matrix_bin} ./cmd/goreveal && "
                        "python3 -m scripts.protected.profile_matrix "
                        f"--goreveal-bin {protected_matrix_bin} "
                        "--output-dir /tmp/protected-matrix "
                        "--json-out /tmp/protected-matrix/profile-matrix.json"
                    ),
                ],
                with_baselines=True,
                env={"GOREVEAL_BIN": protected_matrix_bin},
            ),
            Step(
                cmd=[
                    "bash",
                    "-lc",
                    "cd /workspace && python3 -m unittest scripts.protected.test_profile_matrix",
                ]
            ),
        ],
        "test-snapshots": [
            Step(cmd=["/usr/local/go/bin/go", "test", "./tests/snapshots"]),
        ],
        "snapshot-update": [
            Step(
                cmd=[
                    "/usr/local/go/bin/go",
                    "test",
                    "./tests/snapshots",
                    "-run",
                    "TestAnalyzeFixtureSnapshot",
                    "-update",
                ]
            ),
        ],
        "lint": [
            Step(cmd=["/go/bin/golangci-lint", "run"]),
        ],
        "lint-python": [
            Step(cmd=["ruff", "check", "plugins", "scripts"]),
            Step(cmd=["ty", "check", ".", "--error-on-warning"]),
        ],
        "format-python": [
            Step(cmd=["ruff", "format", "plugins", "scripts"]),
        ],
        "lint-yaml": [
            Step(cmd=["yamllint", ".golangci.yml", ".yamllint.yml", "Taskfile.yml"]),
        ],
        "lint-shell": [
            Step(
                cmd=[
                    "shellcheck",
                    "scripts/baseline/run_gore.sh",
                    "scripts/baseline/run_goresym.sh",
                    "scripts/baseline/run_redress.sh",
                ]
            ),
        ],
        "fuzz": [
            Step(cmd=["/usr/local/go/bin/go", "test", "-fuzz=Fuzz", "-run=^$", "./..."]),
        ],
        "bench": [
            Step(cmd=["/usr/local/go/bin/go", "test", "-bench=.", "-benchmem", "./..."]),
        ],
    }
    if name == "verify":
        return task_steps("fmt") + task_steps("test")
    if name == "lint-scripts":
        return task_steps("lint-python") + task_steps("lint-yaml") + task_steps("lint-shell")
    try:
        return tasks[name]
    except KeyError as exc:
        raise KeyError(f"unknown task: {name}") from exc


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Podman-first GoREveal automation runner.")
    subparsers = parser.add_subparsers(dest="subcommand", required=True)

    build_image = subparsers.add_parser("build-image", help="Build the GoREveal dev image.")
    build_image.add_argument("--image", default=DEFAULT_DEV_IMAGE)

    run_task = subparsers.add_parser("task", help="Run a predefined repository task in Podman.")
    run_task.add_argument("name")

    run_exec = subparsers.add_parser("exec", help="Run an arbitrary command in the dev image.")
    run_exec.add_argument("--with-baselines", action="store_true")
    run_exec.add_argument("command", nargs=argparse.REMAINDER)

    return parser.parse_args(argv)


def main(argv: list[str] | None = None) -> int:
    args = parse_args(sys.argv[1:] if argv is None else argv)
    repo_root = str(Path(__file__).resolve().parents[2])
    client = create_client()

    if args.subcommand == "build-image":
        build_dev_image(client=client, repo_root=repo_root, image=args.image)
        return 0

    if args.subcommand == "task":
        steps = task_steps(args.name)
    else:
        if not args.command:
            raise SystemExit("exec requires a command after '--'")
        command = args.command
        if command[0] == "--":
            command = command[1:]
        steps = [Step(cmd=command, with_baselines=args.with_baselines)]

    for step in steps:
        status = execute_step(client=client, step=step, cwd=repo_root)
        if status != 0:
            return status
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
