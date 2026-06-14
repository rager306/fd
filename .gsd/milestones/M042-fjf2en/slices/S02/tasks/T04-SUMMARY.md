---
id: T04
parent: S02
milestone: M042-fjf2en
key_files:
  - README.md
  - docs/same-host-embedding-service-contract.md
  - docs/fd-v2.md
  - docker-compose.yaml
  - .github/workflows/onnx-packaging.yml
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-06-14T10:46:31.138Z
blocker_discovered: false
---

# T04: Updated operator docs and compose/CI surfaces to TEI-only current posture, with ONNX retained only as historical/future research context.

**Updated operator docs and compose/CI surfaces to TEI-only current posture, with ONNX retained only as historical/future research context.**

## What Happened

Updated README and the same-host embedding service contract to describe TEI as the only current supported runtime and to state that ONNX is disabled as an operator-selectable backend. Added an M042 TEI-only posture note to `docs/fd-v2.md`. Removed outdated compose comments and `--dtype fp16` from the base TEI command so Compose no longer suggests ONNX export as current optimization. Removed the active ONNX packaging GitHub workflow now that `Dockerfile.onnx` and ONNX runtime code have been removed. Historical benchmark/GSD artifacts remain untouched.

## Verification

Verified current docs contain TEI-only posture and no active ONNX operator env instructions such as `ONNX_ARTIFACT_MANIFEST=`, `ONNX_RUNTIME_LIBRARY=`, `ONNX_TOKENIZER_PATH=`, `Dockerfile.onnx`, or ONNX packaging workflow wording. Verified `.github/workflows/onnx-packaging.yml` and `Dockerfile.onnx` are absent. Ran `cd api && go test ./...` successfully after docs/compose changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 - <<'PY'
from pathlib import Path
checks=[]
files=['README.md','docs/same-host-embedding-service-contract.md','docs/fd-v2.md','docker-compose.yaml','docker-compose.override.yaml']
combined='\n'.join(Path(f).read_text() for f in files if Path(f).exists())
for forbidden in ['ONNX_ARTIFACT_MANIFEST=', 'ONNX_RUNTIME_LIBRARY=', 'ONNX_TOKENIZER_PATH=', 'Dockerfile.onnx', 'ONNX Packaging Manual', 'convert model to ONNX using optimum-cli', 'the ONNX backend is opt-in and requires']:
    checks.append((forbidden, forbidden not in combined))
print('tei_only_phrase', 'TEI-only' in combined or 'TEI only' in combined)
for k,v in checks: print(k, v)
PY` | 0 | ✅ pass: current docs have TEI-only posture and no active ONNX operator instructions | 120000ms |
| 2 | `test ! -e .github/workflows/onnx-packaging.yml && echo onnx_workflow_removed=yes; test ! -e Dockerfile.onnx && echo dockerfile_onnx_removed=yes` | 0 | ✅ pass: active ONNX workflow and Dockerfile are removed | 120000ms |
| 3 | `cd api && go test ./...` | 0 | ✅ pass: Go tests still pass after docs/compose cleanup | 180000ms |

## Deviations

ONNX mentions remain where they explicitly state historical/future research or in schema/tests that ensure TEI responses omit ONNX-only fields. Those are not operator instructions.

## Known Issues

TEI's internal ONNX/ORT probing remains in the external TEI image and needs separate TEI startup stabilization research/config; fd docs now distinguish this from fd runtime fallback.

## Files Created/Modified

- `README.md`
- `docs/same-host-embedding-service-contract.md`
- `docs/fd-v2.md`
- `docker-compose.yaml`
- `.github/workflows/onnx-packaging.yml`
