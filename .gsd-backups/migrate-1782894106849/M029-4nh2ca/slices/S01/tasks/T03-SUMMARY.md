---
id: T03
parent: S01
milestone: M029-4nh2ca
key_files:
  - tools/provision_onnx_artifacts.py
key_decisions:
  - S01 accepts GitNexus MEDIUM pre-commit scope because it is confined to `tools/provision_onnx_artifacts.py` and verified through local security probes and provisioning guardrails.
duration: 
verification_result: passed
completed_at: 2026-05-21T04:32:35.204Z
blocker_discovered: false
---

# T03: Verified M029 S01 provisioning hardening guardrails.

**Verified M029 S01 provisioning hardening guardrails.**

## What Happened

Ran S01 guardrails. The provisioning helper compiles, deterministic security probes pass, dry-run JSON output is valid, missing-source behavior still fails as expected, verifier allow-missing still works, default Go tests pass, pinned lint reports 0 issues, actionlint passes, default Docker build passes, binary hygiene passes, port 18000 is clean, and no background processes remain. GitNexus reports expected MEDIUM pre-commit scope confined to the provisioning helper main path.

## Verification

All executable S01 checks passed; GitNexus pre-commit scope recorded as expected.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `python3 -m py_compile tools/provision_onnx_artifacts.py tools/verify_onnx_artifacts.py benchmark.py tools/evaluate_legal_retrieval.py && provisioning dry-run/missing-source/verifier allow-missing` | 0 | ✅ pass — m029_provisioning_guardrails=pass | 6400ms |
| 2 | `cd api && go test ./... -short` | 0 | ✅ pass — Go test: 85 passed in 4 packages | 6300ms |
| 3 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 6200ms |
| 4 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/go-quality.yml .github/workflows/onnx-packaging.yml` | 0 | ✅ pass — no output | 6200ms |
| 5 | `docker build -f api/Dockerfile -t fd-api:m029-default-s01 api` | 0 | ✅ pass — Successfully tagged fd-api:m029-default-s01 | 6100ms |
| 6 | `binary hygiene, port cleanup, bg_shell list` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0; port_18000_clean; no background processes | 0ms |
| 7 | `gitnexus_detect_changes` | 0 | ⚠️ expected pre-commit MEDIUM scope — provisioning helper main path changed | 0ms |

## Deviations

GitNexus pre-commit risk is MEDIUM for the provisioning helper main path, matching the planned security remediation scope. Direct impacts on fetch_source, materialize_source, and parse_args were checked before edits.

## Known Issues

M028 LOW findings remain for future work; final post-commit reindex/detect required.

## Files Created/Modified

- `tools/provision_onnx_artifacts.py`
