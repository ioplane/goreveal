# GoREveal IDA Adapter

<img
  src="https://shieldcn.dev/badge/component-ida--adapter-slate.svg?variant=outline&size=xs"
  alt="component: ida-adapter" height="20">

This adapter is intentionally thin.

Rules:

- it consumes `goreveal export ida` payloads
- it does not recover symbols, types, or strings on its own
- it turns validated export data into a sequence of import actions

Current status:

- pure-Python contract validator and action builder
- container-testable without IDA installed
- ready for later binding to real IDA APIs

Example:

```bash
goreveal export ida ./sample.bin > ida-export.json
python3 plugins/ida/goreveal_ida.py ida-export.json
```
