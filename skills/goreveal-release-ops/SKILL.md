---
name: goreveal-release-ops
description: Manage operational readiness for GoREveal including corpus refresh, regression checks, packaging, and release evidence.
metadata:
  short-description: Release and ops workflow
---

# GoREveal Release and Operations

## Purpose

Use this skill for release prep, CI hardening, and repeatable project operations.

## Required Checks

- architecture docs are still aligned with implementation
- corpus and snapshot expectations are current
- differential tests are green or divergences are documented
- benchmark regressions are reviewed
- release artifacts reflect actual supported capabilities
- Podman verification commands are documented in the current `make`-based workflow
- release claims do not outrun current Sprint 11 metadata surfaces and documented divergence policy

## Scope

- release notes
- CI workflow expectations
- corpus maintenance
- benchmark baselines
- install and operations docs
