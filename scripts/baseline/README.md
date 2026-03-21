# GoREveal Baseline Scripts

This directory will hold wrappers for baseline tool execution and output normalization.

Planned scripts:
- `run_gore.sh`
- `run_redress.sh`
- `run_goresym.sh`
- `run_goresolver.sh`
- `run_gostringungarbler.sh`

Each wrapper should:
- normalize output shape where possible
- avoid mutating baseline repositories
- make differential comparison reproducible

Current state:
- `run_goresym.sh` and `run_redress.sh` now delegate normalization to the shared Python module `scripts.baseline.normalize`
- `run_gore.sh` now also delegates normalization to the shared Python module, while still using a temporary Go helper for extraction because its data path is driven by the `gore` library API rather than plain CLI output
- `scripts.baseline.test_normalize` covers the shared Python normalization logic directly
