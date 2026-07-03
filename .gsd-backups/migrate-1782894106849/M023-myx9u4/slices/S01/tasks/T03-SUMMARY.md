---
id: T03
parent: S01
milestone: M023-myx9u4
key_files:
  - benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt
  - .github/workflows/go-quality.yml
key_decisions:
  - Binary hygiene must distinguish Dockerfile naming (`Dockerfile.onnx`) from actual ONNX model artifacts.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:58:24.368Z
blocker_discovered: false
---

# T03: Validated artifact hygiene, cleaned packaged ONNX runtime, and fixed a CI binary-hygiene false positive.

**Validated artifact hygiene, cleaned packaged ONNX runtime, and fixed a CI binary-hygiene false positive.**

## What Happened

Validated the packaged legal artifact and cleaned up the runtime. The raw legal text leak check reported zero leaks. The artifact exists and the evaluator script compiles. The ONNX container was stopped and port 18000 is clean. During binary hygiene, the existing regex falsely flagged `Dockerfile.onnx`; the workflow hygiene check was corrected and revalidated with actionlint and the refined git-files check.

## Verification

Raw legal leak check, artifact existence, runtime cleanup, and corrected binary hygiene all passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec raw legal text leak check` | 0 | ✅ pass — raw_legal_text_leaks=0 | 55ms |
| 2 | `python3 -m py_compile tools/evaluate_legal_retrieval.py && test -s benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt` | 0 | ✅ pass — m023_legal_artifact_exists=pass | 0ms |
| 3 | `bg_shell kill ebdf984d; docker rm -f fd-onnx-m023-legal; lsof port 18000` | 0 | ✅ pass — port_18000_clean; no background processes | 0ms |
| 4 | `git ls-files refined binary hygiene check` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 0ms |
| 5 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml` | 0 | ✅ pass — actionlint no findings | 4700ms |

## Deviations

The initial binary hygiene command falsely matched tracked `Dockerfile.onnx` as a forbidden model artifact. I fixed the workflow check to exempt `Dockerfile.onnx` while still blocking actual `.onnx` model artifacts, `libtokenizers.a`, and ONNX Runtime shared libraries.

## Known Issues

None for S01. The CI hygiene fix should be included in M023 closure verification.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m023-s01-onnx-docker1024.txt`
- `.github/workflows/go-quality.yml`
