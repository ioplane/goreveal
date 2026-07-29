from __future__ import annotations

import unittest
from pathlib import Path
from unittest import mock

from scripts.protected import profile_matrix


class ProfileMatrixTests(unittest.TestCase):
    def test_platforms_cover_amd64_and_arm64_for_linux_windows_and_darwin(self) -> None:
        actual_platforms = [
            (platform.goos, platform.goarch, platform.binary_name)
            for platform in profile_matrix.PLATFORMS
        ]
        self.assertEqual(
            actual_platforms,
            [
                ("linux", "amd64", "enterprise-sample"),
                ("linux", "arm64", "enterprise-sample"),
                ("windows", "amd64", "enterprise-sample.exe"),
                ("windows", "arm64", "enterprise-sample.exe"),
                ("darwin", "amd64", "enterprise-sample"),
                ("darwin", "arm64", "enterprise-sample"),
            ],
        )

    def test_profiles_cover_first_cross_platform_non_garbled_matrix(self) -> None:
        self.assertEqual(
            [profile.name for profile in profile_matrix.PROFILES],
            [
                "plain",
                "stripped",
                "trimpath",
                "stripped-trimpath",
                "pie",
                "garble",
                "garble-literals-tiny",
            ],
        )

    def test_garble_profiles_are_bounded_to_linux_and_carry_gogarble(self) -> None:
        garble_profiles = {
            profile.name: profile
            for profile in profile_matrix.PROFILES
            if profile.name in {"garble", "garble-literals-tiny"}
        }

        self.assertEqual(set(garble_profiles), {"garble", "garble-literals-tiny"})
        self.assertEqual(garble_profiles["garble"].supported_goos, ("linux",))
        self.assertEqual(garble_profiles["garble"].gogarble, "*")
        self.assertEqual(garble_profiles["garble"].garble_flags, ())
        self.assertEqual(garble_profiles["garble-literals-tiny"].supported_goos, ("linux",))
        self.assertEqual(garble_profiles["garble-literals-tiny"].gogarble, "*")
        self.assertEqual(
            garble_profiles["garble-literals-tiny"].garble_flags,
            ("-literals", "-tiny"),
        )

    def test_profile_record_reads_transfer_and_file_visibility_relevant_surface(self) -> None:
        with mock.patch.object(
            profile_matrix,
            "_baseline_summary",
            return_value={
                "build_info_path": "example.com/protectedfixture",
                "functions": 2,
                "packages": 1,
                "files": 0,
                "relevant_functions": ["main.readLicenseToken"],
            },
        ):
            record = profile_matrix.profile_record(
                profile=profile_matrix.BuildProfile(
                    name="stripped-trimpath",
                    supported_goos=("linux", "windows", "darwin"),
                    trimpath=True,
                    ldflags="-s -w -buildid=",
                ),
                platform=profile_matrix.PlatformTarget(
                    goos="windows",
                    goarch="amd64",
                    binary_name="enterprise-sample.exe",
                ),
                binary_path=Path(
                    "/workspace/.tmp/protected-matrix/windows/amd64/stripped-trimpath/enterprise-sample.exe"
                ),
                analysis_bundle={
                    "analyze": {
                        "input": {"format": "pe"},
                        "build_info": {"path": "example.com/protectedfixture"},
                        "runtime": {
                            "trust_summary": "section_heuristic",
                            "elf_text_section_addr": 0,
                            "elf_text_section_end_inclusive": 0,
                            "elf_pclntab_header_magic": "",
                            "elf_pclntab_header_magic_kind": "",
                            "elf_pclntab_function_count_hint": 0,
                            "elf_pclntab_file_count_hint": 0,
                            "elf_functab_first_pc_offset_hint": 0,
                            "elf_functab_last_pc_offset_hint": 0,
                            "elf_functab_pc_offsets_monotonic": False,
                            "elf_functab_first_pc_addr_hint": 0,
                            "elf_functab_last_pc_addr_hint": 0,
                            "elf_functab_pc_addr_hints_within_text": False,
                            "elf_functab_pc_addr_sample": [],
                            "elf_functab_pc_addr_sample_all_within_text": False,
                            "elf_function_foothold": "",
                            "elf_function_foothold_count_hint": 0,
                            "elf_function_foothold_text_source": "",
                            "elf_function_foothold_start_addr": 0,
                            "elf_function_foothold_end_addr": 0,
                            "elf_function_recovery_blocker": "",
                            "moduledata_pcheader_matches_gopclntab": False,
                            "moduledata_funcnametab_within_gopclntab": False,
                            "moduledata_pclntable_within_gopclntab": False,
                        },
                        "functions": [
                            {"name": "main.readLicenseToken"},
                            {"name": "main.auditFeatureGate"},
                        ],
                        "packages": [{}],
                        "source_tree": {"files": ["main.go"], "pathless_file_evidence": True},
                    },
                    "peel": {
                        "functions": [
                            {"name": "main.readLicenseToken"},
                            {"name": "main.auditFeatureGate"},
                        ],
                        "packages": [{"primary_classification": "user"}],
                    },
                },
            )

        self.assertEqual(record["profile"], "stripped-trimpath")
        self.assertEqual(record["target"], "windows-amd64")
        self.assertEqual(record["format"], "pe")
        self.assertEqual(record["goreveal"]["runtime_trust_summary"], "section_heuristic")
        self.assertEqual(record["goreveal"]["elf_text_section_addr"], 0)
        self.assertEqual(record["goreveal"]["elf_text_section_end_inclusive"], 0)
        self.assertEqual(record["goreveal"]["elf_pclntab_header_magic"], "")
        self.assertEqual(record["goreveal"]["elf_pclntab_header_magic_kind"], "")
        self.assertEqual(record["goreveal"]["elf_pclntab_function_count_hint"], 0)
        self.assertEqual(record["goreveal"]["elf_pclntab_file_count_hint"], 0)
        self.assertEqual(record["goreveal"]["elf_functab_first_pc_offset_hint"], 0)
        self.assertEqual(record["goreveal"]["elf_functab_last_pc_offset_hint"], 0)
        self.assertFalse(record["goreveal"]["elf_functab_pc_offsets_monotonic"])
        self.assertEqual(record["goreveal"]["elf_functab_first_pc_addr_hint"], 0)
        self.assertEqual(record["goreveal"]["elf_functab_last_pc_addr_hint"], 0)
        self.assertFalse(record["goreveal"]["elf_functab_pc_addr_hints_within_text"])
        self.assertEqual(record["goreveal"]["elf_functab_pc_addr_sample_count"], 0)
        self.assertEqual(record["goreveal"]["elf_functab_pc_addr_sample_first"], 0)
        self.assertFalse(record["goreveal"]["elf_functab_pc_addr_sample_all_within_text"])
        self.assertEqual(record["goreveal"]["elf_function_foothold"], "")
        self.assertEqual(record["goreveal"]["elf_function_foothold_count_hint"], 0)
        self.assertEqual(record["goreveal"]["elf_function_foothold_text_source"], "")
        self.assertEqual(record["goreveal"]["elf_function_foothold_start_addr"], 0)
        self.assertEqual(record["goreveal"]["elf_function_foothold_end_addr"], 0)
        self.assertEqual(record["goreveal"]["elf_function_recovery_blocker"], "")
        self.assertFalse(record["goreveal"]["moduledata_pcheader_matches_gopclntab"])
        self.assertFalse(record["goreveal"]["moduledata_funcnametab_within_gopclntab"])
        self.assertFalse(record["goreveal"]["moduledata_pclntable_within_gopclntab"])
        self.assertEqual(record["goreveal"]["functions"], 2)
        self.assertEqual(record["goreveal"]["files"], 1)
        self.assertTrue(record["goreveal"]["pathless_file_evidence"])
        self.assertEqual(
            record["goreveal"]["relevant_functions"],
            ["main.readLicenseToken", "main.auditFeatureGate"],
        )

    def test_tool_function_presence_filters_relevant_main_functions(self) -> None:
        present = profile_matrix._tool_function_presence(
            [
                "main.main",
                "main.runEnterpriseReport",
                "main.readLicenseToken",
                "internal/licensegate.VerifyLicenseToken",
            ]
        )

        self.assertEqual(
            present,
            ["main.readLicenseToken", "main.runEnterpriseReport"],
        )

    @mock.patch("scripts.protected.profile_matrix.subprocess.run")
    def test_build_go_binary_prefixes_toolchain_path_for_garble(self, run_mock: mock.Mock) -> None:
        profile_matrix.build_go_binary(
            "/usr/local/go/bin/go",
            "/go/bin/garble",
            Path("/tmp/source"),
            Path("/tmp/out.bin"),
            profile_matrix.BuildProfile(
                name="garble",
                supported_goos=("linux",),
                gogarble="*",
            ),
            profile_matrix.PlatformTarget(
                goos="linux",
                goarch="amd64",
                binary_name="enterprise-sample",
            ),
        )

        _, kwargs = run_mock.call_args
        env = kwargs["env"]
        self.assertTrue(env["PATH"].startswith("/usr/local/go/bin:/go/bin:"))
        self.assertEqual(env["GOGARBLE"], "*")

    @mock.patch("scripts.protected.profile_matrix.subprocess.run")
    def test_ensure_garble_binary_builds_local_source_when_present(
        self, run_mock: mock.Mock
    ) -> None:
        with mock.patch.object(Path, "exists", return_value=True):
            resolved = profile_matrix.ensure_garble_binary(
                go_bin="/usr/local/go/bin/go",
                garble_bin="/go/bin/garble",
                garble_source=Path("/repos/garble"),
                tools_dir=Path("/tmp/tools"),
            )

        _, kwargs = run_mock.call_args
        self.assertEqual(kwargs["cwd"], Path("/repos/garble"))
        self.assertEqual(kwargs["env"]["GOWORK"], "off")
        self.assertEqual(resolved, "/tmp/tools/garble")

    def test_ensure_garble_binary_falls_back_to_installed_binary_when_source_missing(self) -> None:
        resolved = profile_matrix.ensure_garble_binary(
            go_bin="/usr/local/go/bin/go",
            garble_bin="/go/bin/garble",
            garble_source=Path("/repos/garble-does-not-exist"),
            tools_dir=Path("/tmp/tools"),
        )

        self.assertEqual(resolved, "/go/bin/garble")


if __name__ == "__main__":
    unittest.main()
