import os
import subprocess
import unittest
from pathlib import Path

from plugins.ida.goreveal_ida import apply_actions, build_actions, load_export


class FakeIDAAPI:
    def __init__(self):
        self.calls = []

    def ensure_function(self, start, end):
        self.calls.append(("ensure_function", start, end))

    def rename_function(self, start, name):
        self.calls.append(("rename_function", start, name))

    def set_comment(self, ea, comment):
        self.calls.append(("set_comment", ea, comment))

    def ensure_folder(self, name):
        self.calls.append(("ensure_folder", name))


class GoRevealIDAAdapterTests(unittest.TestCase):
    def test_build_actions_prefers_refined_names(self):
        payload = {
            "contract": "goreveal.export.ida/v1",
            "functions": [
                {
                    "name": "main.main",
                    "refined_name": "main.main",
                    "entry": 4096,
                    "end": 4352,
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
            "packages": [
                {
                    "name": "main",
                    "function_count": 1,
                }
            ],
            "source_tree": {
                "root": "example.com/sample",
                "files": ["main.go"],
            },
        }

        actions = build_actions(payload)

        self.assertEqual(actions[0]["kind"], "function")
        self.assertEqual(actions[0]["name"], "main.main")
        self.assertEqual(actions[1]["kind"], "type")
        self.assertEqual(actions[2]["value"], "hello refined")
        self.assertEqual(actions[2]["address"], 8192)
        self.assertEqual(actions[3]["kind"], "package")
        self.assertEqual(actions[4]["kind"], "source_file")

    def test_load_export_rejects_wrong_contract(self):
        with self.assertRaises(ValueError):
            load_export({"contract": "goreveal.export.ghidra/v1"})

    def test_apply_actions_uses_fake_api(self):
        actions = [
            {"kind": "function", "start": 4096, "end": 4352, "name": "main.main"},
            {"kind": "package", "name": "main", "function_count": 2},
            {"kind": "string", "address": 8192, "offset": 16, "value": "hello"},
        ]
        api = FakeIDAAPI()

        apply_actions(api, actions)

        self.assertIn(("ensure_function", 4096, 4352), api.calls)
        self.assertIn(("rename_function", 4096, "main.main"), api.calls)
        self.assertIn(("ensure_folder", "main"), api.calls)
        self.assertIn(("set_comment", 8192, "hello"), api.calls)

    def test_fixture_export_builds_actionable_payload(self):
        repo_root = Path(__file__).resolve().parents[2]
        fixture = repo_root / "corpus" / "fixtures" / "go-elf-buildinfo-linux-amd64" / "fixture.bin"
        goreveal_bin = os.environ.get("GOREVEAL_BIN")
        cmd = [goreveal_bin, "export", "ida", str(fixture)] if goreveal_bin else [
            "/usr/local/go/bin/go",
            "run",
            "./cmd/goreveal",
            "export",
            "ida",
            str(fixture),
        ]

        proc = subprocess.run(
            cmd,
            cwd=repo_root,
            check=True,
            capture_output=True,
            text=True,
        )

        actions = build_actions(load_export(proc.stdout))

        self.assertTrue(any(action["kind"] == "function" and action["name"] == "main.main" for action in actions))
        self.assertTrue(any(action["kind"] == "package" and action["name"] == "main" for action in actions))
        self.assertTrue(any(action["kind"] == "source_file" and action["name"] == "main.go" for action in actions))


if __name__ == "__main__":
    unittest.main()
