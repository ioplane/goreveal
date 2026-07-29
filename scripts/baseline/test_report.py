import unittest

from scripts.baseline.report import build_fixture_report


class BuildFixtureReportTest(unittest.TestCase):
    def test_builds_match_skip_and_divergence_summary(self) -> None:
        analysis = {
            "build_info": {
                "go_version": "go1.26.1",
                "path": "example.com/gorevealfixture",
            },
            "packages": [
                {"name": "main"},
            ],
            "source_tree": {
                "root": "example.com/gorevealfixture",
                "files": ["main.go"],
            },
        }
        goresym = {
            "build_info": {
                "go_version": "go1.26.1",
                "path": "example.com/gorevealfixture",
            },
            "files": ["example.com/gorevealfixture/main.go"],
            "functions": ["main.main", "main.helperAdd", "main.helperBanner"],
        }
        redress = {
            "build_info": {
                "path": "example.com/gorevealfixture",
            },
            "packages": ["main"],
            "source_files": ["main.go"],
            "functions": ["main.main", "main.helperBanner"],
        }
        gore = {
            "build_info": {
                "go_version": "go1.26.1",
                "path": "wrong.example/module",
            },
            "packages": ["main"],
            "source_files": ["main.go"],
            "functions": ["main.main", "main.helperBanner"],
            "types": ["fixtureCounter"],
        }

        report = build_fixture_report(
            fixture="go-elf-buildinfo-linux-amd64",
            analysis=analysis,
            goresym=goresym,
            redress=redress,
            gore=gore,
        )

        self.assertEqual(report["fixture"], "go-elf-buildinfo-linux-amd64")
        self.assertEqual(report["summary"]["diverged"], 1)
        self.assertGreaterEqual(report["summary"]["matched"], 1)
        self.assertGreaterEqual(report["summary"]["skipped"], 1)

        statuses = {
            (entry["baseline"], entry["subject"]): entry["status"] for entry in report["checks"]
        }
        self.assertEqual(statuses[("gore", "build_info.path")], "diverged")
        self.assertEqual(statuses[("goresym", "build_info.go_version")], "matched")
        self.assertEqual(statuses[("redress", "build_info.go_version")], "skipped")


if __name__ == "__main__":
    unittest.main()
