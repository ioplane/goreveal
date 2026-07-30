import os
import subprocess
import unittest
from pathlib import Path
from typing import Any

from plugins.ghidra.goreveal_ghidra import apply_actions, build_actions, load_export


class FakeGhidraAPI:
    def __init__(self) -> None:
        self.calls: list[tuple[str, Any, Any, Any] | tuple[str, Any, Any] | tuple[str, Any]] = []

    def set_program_context(self, path: str, module_path: str, go_version: str) -> None:
        self.calls.append(("set_program_context", path, module_path, go_version))

    def ensure_function(self, start: int, end: int) -> None:
        self.calls.append(("ensure_function", start, end))

    def rename_symbol(self, address: int, name: str) -> None:
        self.calls.append(("rename_symbol", address, name))

    def ensure_namespace(self, name: str) -> None:
        self.calls.append(("ensure_namespace", name))

    def set_comment(self, ea: int, comment: str) -> None:
        self.calls.append(("set_comment", ea, comment))


class GoRevealGhidraAdapterTests(unittest.TestCase):
    def test_build_actions_prefers_refined_names(self) -> None:
        payload = {
            "contract": "goreveal.export.ghidra/v1",
            "program": {
                "path": "/tmp/sample.bin",
                "module_path": "example.com/sample",
                "go_version": "go1.26.1",
            },
            "symbols": [
                {
                    "name": "main.main",
                    "refined_name": "main.main",
                    "address": 4096,
                    "end": 4352,
                }
            ],
            "packages": [
                {
                    "name": "main",
                    "function_count": 1,
                }
            ],
            "types": [
                {
                    "name": "main.counter",
                    "refined_name": "main.counter",
                    "kind": "struct",
                }
            ],
            "strings": [
                {
                    "value": "hello",
                    "refined_value": "hello refined",
                    "address": 8192,
                    "offset": 16,
                    "region": ".rodata",
                }
            ],
            "source_tree": {
                "root": "example.com/sample",
                "files": ["main.go"],
            },
        }

        actions = build_actions(payload)

        self.assertEqual(actions[0]["kind"], "program")
        self.assertEqual(actions[1]["kind"], "symbol")
        self.assertEqual(actions[1]["name"], "main.main")
        self.assertEqual(actions[2]["kind"], "package")
        self.assertEqual(actions[3]["kind"], "type")
        self.assertEqual(actions[4]["value"], "hello refined")
        self.assertEqual(actions[4]["address"], 8192)
        self.assertEqual(actions[5]["kind"], "source_file")

    def test_load_export_rejects_wrong_contract(self) -> None:
        with self.assertRaises(ValueError):
            load_export({"contract": "goreveal.export.ida/v1"})

    def test_load_export_preserves_runtime_surface(self) -> None:
        payload = load_export(
            {
                "contract": "goreveal.export.ghidra/v1",
                "program": {
                    "path": "/tmp/sample.bin",
                    "module_path": "example.com/sample",
                    "go_version": "go1.26.1",
                },
                "runtime": {
                    "trust_summary": "go_module_fallback",
                    "elf_function_foothold": "address_only",
                    "elf_function_foothold_count_hint": 2083,
                    "elf_function_foothold_text_source": "elf_text_section",
                    "elf_function_foothold_start_addr": 0x11000,
                    "elf_function_foothold_end_addr": 0xB55D1,
                },
                "symbols": [],
            }
        )

        runtime = payload.get("runtime") or {}
        self.assertEqual(runtime.get("trust_summary"), "go_module_fallback")
        self.assertEqual(runtime.get("elf_function_foothold"), "address_only")
        self.assertEqual(runtime.get("elf_function_foothold_count_hint"), 2083)
        self.assertEqual(runtime.get("elf_function_foothold_text_source"), "elf_text_section")
        self.assertEqual(runtime.get("elf_function_foothold_start_addr"), 0x11000)
        self.assertEqual(runtime.get("elf_function_foothold_end_addr"), 0xB55D1)

    def test_apply_actions_uses_fake_api(self) -> None:
        actions = [
            {
                "kind": "program",
                "path": "/tmp/sample.bin",
                "module_path": "example.com/sample",
                "go_version": "go1.26.1",
            },
            {"kind": "symbol", "address": 4096, "end": 4352, "name": "main.main"},
            {"kind": "package", "name": "main", "function_count": 2},
            {"kind": "string", "address": 8192, "offset": 16, "value": "hello"},
        ]
        api = FakeGhidraAPI()

        apply_actions(api, actions)

        self.assertIn(
            ("set_program_context", "/tmp/sample.bin", "example.com/sample", "go1.26.1"),
            api.calls,
        )
        self.assertIn(("ensure_function", 4096, 4352), api.calls)
        self.assertIn(("rename_symbol", 4096, "main.main"), api.calls)
        self.assertIn(("ensure_namespace", "main"), api.calls)
        self.assertIn(("set_comment", 8192, "hello"), api.calls)

    def test_fixture_export_builds_actionable_payload(self) -> None:
        repo_root = Path(__file__).resolve().parents[2]
        fixture = repo_root / "corpus" / "fixtures" / "go-elf-buildinfo-linux-amd64" / "fixture.bin"
        goreveal_bin = os.environ.get("GOREVEAL_BIN")
        cmd = (
            [goreveal_bin, "export", "ghidra", str(fixture)]
            if goreveal_bin
            else [
                "/usr/local/go/bin/go",
                "run",
                "./cmd/goreveal",
                "export",
                "ghidra",
                str(fixture),
            ]
        )

        proc = subprocess.run(
            cmd,
            cwd=repo_root,
            check=True,
            capture_output=True,
            text=True,
        )

        payload = load_export(proc.stdout)
        actions = build_actions(payload)

        runtime = payload.get("runtime") or {}
        self.assertEqual(runtime.get("trust_summary"), "symbol_backed")
        self.assertIn(runtime.get("elf_function_foothold"), {"", None, "address_only"})

        self.assertTrue(
            any(
                action["kind"] == "program"
                and action["module_path"] == "example.com/gorevealfixture"
                for action in actions
            )
        )
        self.assertTrue(
            any(action["kind"] == "symbol" and action["name"] == "main.main" for action in actions)
        )
        self.assertTrue(
            any(action["kind"] == "package" and action["name"] == "main" for action in actions)
        )
        self.assertTrue(
            any(
                action["kind"] == "source_file" and action["name"] == "main.go"
                for action in actions
            )
        )


if __name__ == "__main__":
    unittest.main()
