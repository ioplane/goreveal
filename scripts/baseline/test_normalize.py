import unittest

from scripts.baseline.normalize import (
    normalize_gore,
    normalize_goresym,
    normalize_redress,
)


class NormalizeGoReSymTest(unittest.TestCase):
    def test_extracts_build_info_files_and_user_functions(self) -> None:
        raw = {
            "BuildInfo": {
                "GoVersion": "go1.26.1",
                "Path": "example.com/gorevealfixture",
            },
            "Files": [
                "example.com/gorevealfixture/main.go",
                "runtime/proc.go",
            ],
            "UserFunctions": [
                {"FullName": "main.main"},
                {"FullName": "main.helperAdd"},
                {"FullName": "main.main"},
                {"FullName": ""},
            ],
        }

        normalized = normalize_goresym(raw)

        self.assertEqual(
            normalized,
            {
                "build_info": {
                    "go_version": "go1.26.1",
                    "path": "example.com/gorevealfixture",
                },
                "files": [
                    "example.com/gorevealfixture/main.go",
                    "runtime/proc.go",
                ],
                "functions": [
                    "main.helperAdd",
                    "main.main",
                ],
            },
        )


class NormalizeRedressTest(unittest.TestCase):
    def test_extracts_packages_source_files_and_functions(self) -> None:
        packages_output = """
Packages:
Name  Version
----  -------
main
main
"""
        source_output = """
Package main: example.com/gorevealfixture
File: main.go
\thelperAdd Lines: 10 to 12 (2)
\thelperBanner Lines: 14 to 16 (2)
\tmain Lines: 18 to 22 (4)
"""
        gomod_output = """
Type  Name                         Version  Replaced by  Hash
----  ----                         -------  -----------  ----
main  example.com/gorevealfixture
"""

        normalized = normalize_redress(packages_output, source_output, gomod_output)

        self.assertEqual(
            normalized,
            {
                "build_info": {
                    "path": "example.com/gorevealfixture",
                },
                "packages": ["main"],
                "source_files": ["main.go"],
                "functions": [
                    "main.helperAdd",
                    "main.helperBanner",
                    "main.main",
                ],
            },
        )


class NormalizeGoreTest(unittest.TestCase):
    def test_extracts_and_sorts_core_fields(self) -> None:
        raw = {
            "build_info": {
                "go_version": "go1.26.1",
                "path": "example.com/gorevealfixture",
            },
            "packages": ["main", "main"],
            "source_files": ["main.go", "main.go"],
            "functions": ["main.main", "main.helperBanner", "main.main"],
            "types": ["fixtureGreeterImpl", "fixtureCounter", "fixtureCounter"],
        }

        normalized = normalize_gore(raw)

        self.assertEqual(
            normalized,
            {
                "build_info": {
                    "go_version": "go1.26.1",
                    "path": "example.com/gorevealfixture",
                },
                "packages": ["main"],
                "source_files": ["main.go"],
                "functions": [
                    "main.helperBanner",
                    "main.main",
                ],
                "types": [
                    "fixtureCounter",
                    "fixtureGreeterImpl",
                ],
            },
        )


if __name__ == "__main__":
    unittest.main()
