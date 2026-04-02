import tempfile
import unittest
from pathlib import Path

from scripts.dev.podman_runner import (
    DEFAULT_BASELINES_ROOT,
    DEFAULT_DEV_IMAGE,
    DEFAULT_DEV_WORKDIR,
    Step,
    default_base_url,
    run_kwargs_for_step,
    task_steps,
)


class DefaultBaseURLTests(unittest.TestCase):
    def test_prefers_explicit_environment_variable(self) -> None:
        env = {"PODMAN_BASE_URL": "unix:///tmp/custom/podman.sock"}

        self.assertEqual(default_base_url(env=env, uid=9999), "unix:///tmp/custom/podman.sock")

    def test_prefers_container_host_environment_variable(self) -> None:
        env = {
            "CONTAINER_HOST": "unix:///run/podman/podman.sock",
            "XDG_RUNTIME_DIR": "/tmp/runtime-dir",
        }

        self.assertEqual(default_base_url(env=env, uid=9999), "unix:///run/podman/podman.sock")

    def test_prefers_docker_host_environment_variable(self) -> None:
        env = {
            "DOCKER_HOST": "unix:///run/podman/podman.sock",
            "XDG_RUNTIME_DIR": "/tmp/runtime-dir",
        }

        self.assertEqual(default_base_url(env=env, uid=9999), "unix:///run/podman/podman.sock")

    def test_uses_xdg_runtime_dir_when_present(self) -> None:
        env = {"XDG_RUNTIME_DIR": "/tmp/runtime-dir"}

        self.assertEqual(
            default_base_url(env=env, uid=1000),
            "unix:///tmp/runtime-dir/podman/podman.sock",
        )

    def test_falls_back_to_run_user_uid(self) -> None:
        self.assertEqual(
            default_base_url(env={}, uid=1234),
            "unix:///run/user/1234/podman/podman.sock",
        )


class RunKwargsForStepTests(unittest.TestCase):
    def test_mounts_workspace_and_sets_working_directory(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            kwargs = run_kwargs_for_step(
                Step(cmd=["/usr/local/go/bin/go", "test", "./..."]),
                cwd=tmpdir,
            )

        self.assertEqual(kwargs["image"], DEFAULT_DEV_IMAGE)
        self.assertEqual(kwargs["working_dir"], DEFAULT_DEV_WORKDIR)
        self.assertEqual(
            kwargs["volumes"],
            {
                str(Path(tmpdir).resolve()): {
                    "bind": DEFAULT_DEV_WORKDIR,
                    "mode": "Z",
                }
            },
        )

    def test_adds_baseline_mount_and_env_when_requested(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            kwargs = run_kwargs_for_step(
                Step(
                    cmd=["/usr/local/go/bin/go", "test", "./tests/differential"],
                    with_baselines=True,
                ),
                cwd=tmpdir,
            )

        self.assertEqual(
            kwargs["volumes"][DEFAULT_BASELINES_ROOT],
            {"bind": "/repos", "mode": "Z"},
        )
        self.assertEqual(kwargs["environment"]["GOREVEAL_BASELINES_ROOT"], "/repos")

    def test_merges_custom_environment(self) -> None:
        with tempfile.TemporaryDirectory() as tmpdir:
            kwargs = run_kwargs_for_step(
                Step(
                    cmd=["bash", "-lc", "echo ok"],
                    env={"GOREVEAL_BIN": f"{DEFAULT_DEV_WORKDIR}/.tmp/goreveal"},
                ),
                cwd=tmpdir,
            )

        self.assertEqual(kwargs["environment"]["GOREVEAL_BIN"], "/workspace/.tmp/goreveal")


class TaskStepsTests(unittest.TestCase):
    def test_build_dev_bin_is_single_container_step(self) -> None:
        steps = task_steps("build-dev-bin")

        self.assertEqual(len(steps), 1)
        self.assertEqual(
            steps[0].cmd,
            [
                "bash",
                "-lc",
                "mkdir -p .tmp && /usr/local/go/bin/go build -o .tmp/goreveal ./cmd/goreveal",
            ],
        )

    def test_test_task_builds_binary_before_running_regression_steps(self) -> None:
        steps = task_steps("test")

        self.assertGreaterEqual(len(steps), 5)
        self.assertEqual(steps[0].cmd, task_steps("build-dev-bin")[0].cmd)
        self.assertTrue(steps[1].with_baselines)
        self.assertIn("GOREVEAL_BIN", steps[4].env)
        self.assertIn("GOREVEAL_BIN", steps[5].env)

    def test_test_differential_report_uses_built_workspace_binary(self) -> None:
        steps = task_steps("test-differential-report")

        self.assertEqual(steps[0].cmd, task_steps("build-dev-bin")[0].cmd)
        self.assertTrue(steps[1].with_baselines)
        self.assertEqual(steps[1].env["GOREVEAL_BIN"], f"{DEFAULT_DEV_WORKDIR}/.tmp/goreveal")

    def test_lint_scripts_expands_into_python_yaml_and_shell_steps(self) -> None:
        steps = task_steps("lint-scripts")

        assert len(steps) == 4
        assert steps[0].cmd == ["ruff", "check", "plugins", "scripts"]
        assert steps[1].cmd == ["ty", "check", ".", "--error-on-warning"]
        assert steps[2].cmd == ["yamllint", ".golangci.yml", ".yamllint.yml", "Taskfile.yml"]
        assert steps[3].cmd[0] == "shellcheck"

    def test_protected_matrix_builds_fresh_binary_in_same_container_step(self) -> None:
        steps = task_steps("protected-matrix")
        expected_cmd = (
            "cd /workspace && "
            "mkdir -p .tmp/protected /tmp/protected-matrix && "
            "/usr/local/go/bin/go build -o /workspace/.tmp/protected/goreveal "
            "./cmd/goreveal && "
            "python3 -m scripts.protected.profile_matrix "
            "--goreveal-bin /workspace/.tmp/protected/goreveal "
            "--output-dir /tmp/protected-matrix "
            "--json-out /tmp/protected-matrix/profile-matrix.json"
        )

        self.assertEqual(len(steps), 2)
        self.assertTrue(steps[0].with_baselines)
        self.assertEqual(steps[0].env["GOREVEAL_BIN"], "/workspace/.tmp/protected/goreveal")
        self.assertEqual(
            steps[0].cmd,
            [
                "bash",
                "-lc",
                expected_cmd,
            ],
        )
        self.assertEqual(
            steps[1].cmd,
            [
                "bash",
                "-lc",
                "cd /workspace && python3 -m unittest scripts.protected.test_profile_matrix",
            ],
        )


if __name__ == "__main__":
    unittest.main()
