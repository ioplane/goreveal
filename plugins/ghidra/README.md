# GoREveal Ghidra Adapter

This adapter is intentionally thin.

Rules:
- it consumes `goreveal export ghidra` payloads
- it does not recover symbols, types, or strings on its own
- it turns validated export data into a sequence of import actions

Current status:
- pure-Python contract validator and action builder
- container-testable without Ghidra installed
- ready for later binding to real Ghidra APIs

Example:
```bash
goreveal export ghidra ./sample.bin > ghidra-export.json
python3 plugins/ghidra/goreveal_ghidra.py ghidra-export.json
```
