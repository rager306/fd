---
id: T04
parent: S02
milestone: M038-pmw50e
key_files:
  - benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt
  - benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt
  - benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt
  - benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:48:02.477Z
blocker_discovered: false
---

# T04: Ran final M038 guardrails successfully.

**Ran final M038 guardrails successfully.**

## What Happened

Ran final M038 guardrails after Go target-runtime smoke, legal, and performance evidence. Manifest/tooling checks, actionlint, Go tests, lint, tagged tokenizer/ONNX tests, Docker default build, evidence leak checks, tracked binary hygiene, background/port checks, and GitNexus detect all passed.

## Verification

All final M038 checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `json.tool + py_compile + provisioning dry-run + verify_onnx_artifacts + verify_onnx_export_contract` | 0 | ✅ pass — manifest valid and provisioning/export checks passed | 25900ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 25900ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 25900ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 25800ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 13900ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 13800ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m038-default-final api` | 0 | ✅ pass — image fd-api:m038-default-final built | 13700ms |
| 8 | `M038 evidence leak checks` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 25700ms |
| 9 | `tracked binary hygiene` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 13700ms |
| 10 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk, no affected processes | 0ms |
| 11 | `bg_shell list and lsof port 18000` | 0 | ✅ pass — no background processes; port_18000_clean | 0ms |

## Deviations

None.

## Known Issues

Packaged Docker ONNX legal/performance reruns and hosted workflow proof remain future gates. Benchmark Redis L2 restart subchecks were skipped due to bg_shell-managed server.

## Files Created/Modified

- `benchmark-results/fd-onnx-go-runtime-smoke-m038-s01.txt`
- `benchmark-results/fd-legal-retrieval-m038-go-onnx-target-runtime.txt`
- `benchmark-results/fd-benchmark-m038-go-onnx-target-runtime.txt`
- `benchmark-results/fd-onnx-go-target-runtime-acceptance-m038-s02.txt`
