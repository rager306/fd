---
id: T01
parent: S02
milestone: M022-i079tk
key_files:
  - .github/workflows/go-quality.yml
key_decisions:
  - CI will run `tools/verify_onnx_artifacts.py --allow-missing` so metadata/schema/default-production contract is checked without requiring binary artifacts.
  - CI will fail if `.onnx`, `libtokenizers.a`, or `libonnxruntime.so` are tracked.
  - Workflow triggers now include ONNX packaging docs/tooling paths.
duration: 
verification_result: passed
completed_at: 2026-05-20T10:45:11.174Z
blocker_discovered: false
---

# T01: Added CI-safe ONNX artifact metadata and binary hygiene checks.

**Added CI-safe ONNX artifact metadata and binary hygiene checks.**

## What Happened

Updated Go Quality workflow with CI-safe ONNX artifact contract checks. The workflow now triggers on ONNX packaging docs/tool changes, verifies manifest metadata in allow-missing mode, and enforces binary hygiene. `actionlint` passed for the workflow file, the local allow-missing verifier passed, and the binary hygiene command passed.

## Verification

Workflow validation and local CI-equivalent artifact checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml` | 0 | ✅ pass — actionlint produced no findings | 5400ms |
| 2 | `python3 -m py_compile tools/verify_onnx_artifacts.py && python3 tools/verify_onnx_artifacts.py --allow-missing ...` | 0 | ✅ pass — ci_allow_missing_verifier=pass | 0ms |
| 3 | `git ls-files | grep -E '(\.onnx$|libtokenizers\.a$|libonnxruntime\.so)'` | 0 | ✅ pass — ci_binary_hygiene=pass | 0ms |

## Deviations

Ruby was unavailable for generic YAML parsing, so I used `actionlint` via `go run` for stronger workflow validation.

## Known Issues

Full ONNX Docker image build is intentionally not added to hosted CI yet because the binary artifact provisioning store/cache is not implemented.

## Files Created/Modified

- `.github/workflows/go-quality.yml`
