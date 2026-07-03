---
id: T01
parent: S01
milestone: M021-4t2wpt
key_files:
  - api/Dockerfile
  - .github/workflows/go-quality.yml
  - .gitignore
  - docs/onnx-artifacts/user-bge-m3-dense-fp32.json
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - Use a script+README artifact provisioning contract in M021, not Dockerfile/CI changes yet.
  - Keep default Dockerfile and GitHub Actions independent of ONNX/native artifacts; ONNX packaging automation remains the next gate after local artifact verification contract.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:13:19.311Z
blocker_discovered: false
---

# T01: Inspected packaging boundaries and chose a script+README contract without changing default Docker/CI.

**Inspected packaging boundaries and chose a script+README contract without changing default Docker/CI.**

## What Happened

Inspected the default Dockerfile, GitHub Actions workflow, `.gitignore`, and ONNX/native manifests. The default Dockerfile builds a CGO-disabled static binary and should remain independent from native tokenizer artifacts. The CI workflow runs default Go tests/lint only. `.gitignore` excludes `.gsd/runtime/`, `*.a`, and model directories. Therefore the minimal safe M021 contract is a verifier script plus README that documents artifact provisioning and checksum requirements without changing default Docker/CI behavior.

## Verification

Relevant files were read and new artifact contract paths are available.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read api/Dockerfile .github/workflows/go-quality.yml .gitignore manifests` | 0 | ✅ pass — packaging boundaries inspected | 0ms |
| 2 | `test ! -e tools/verify_onnx_artifacts.py && test ! -e docs/onnx-artifacts/README.md` | 0 | ✅ pass — new_artifact_contract_paths_ok | 0ms |

## Deviations

None.

## Known Issues

Current CI does not provision local ignored ONNX/native tokenizer artifacts, by design. A future packaging milestone must add CI artifact provisioning or a controlled skip/contract check.

## Files Created/Modified

- `api/Dockerfile`
- `.github/workflows/go-quality.yml`
- `.gitignore`
- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
