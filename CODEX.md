# Codex Notes for GoREveal

Read `AGENTS.md` first.

Use this file as an implementation overlay.

Focus areas:
- narrow vertical slices
- test-backed changes
- schema-safe implementation
- benchmark-backed performance work

When working here:
- implement recovery logic with matching schema mapping and evidence
- do not introduce interfaces before a second real implementation exists
- do not optimize before a hotspot is measured
- keep scalar and optimized paths behaviorally identical


Container rule:
- do not rely on host Go or host linters
- run build, test, lint, fuzz, and bench actions inside the Podman dev container
