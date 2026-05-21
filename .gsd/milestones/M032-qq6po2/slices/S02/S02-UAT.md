# S02: Reproducibility strategy documentation and closure — UAT

**Milestone:** M032-qq6po2
**Written:** 2026-05-21T07:02:39.770Z

# UAT — M032 S02

A future operator can:

1. run `python3 tools/verify_onnx_export_contract.py`;
2. see whether the current local `.onnx` artifact matches manifest/provenance/export metadata;
3. read `docs/onnx-artifacts/PROVISIONING.md` to understand that this is not regenerated export proof;
4. choose the next gate: exact-binary immutable hosting or full reproducible-export workflow.

The docs and outcome explicitly preserve TEI as production/default and ONNX as opt-in experimental.

