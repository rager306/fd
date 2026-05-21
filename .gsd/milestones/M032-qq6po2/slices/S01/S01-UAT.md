# S01: Local export contract verifier — UAT

**Milestone:** M032-qq6po2
**Written:** 2026-05-21T06:57:45.005Z

# UAT — M032 S01

A future operator can run:

```bash
python3 tools/verify_onnx_export_contract.py
```

and receive structured JSON showing whether the current local ONNX artifact matches the tracked manifest/provenance/export metadata. The output explicitly states that this is existing-artifact contract verification, not regenerated export proof.

