---
id: T04
parent: S02
milestone: M039-aexhf5
key_files:
  - benchmark-results/fd-onnx-docker-smoke-m039-s01.txt
  - benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt
  - benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt
  - benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T11:32:09.733Z
blocker_discovered: false
---

# T04: Ran final M039 guardrails successfully.

**Ran final M039 guardrails successfully.**

## What Happened

Ran final M039 guardrails after packaged smoke/legal/performance evidence. Manifest/tooling checks, actionlint, Go tests, lint, tagged tokenizer/ONNX tests, Docker default build, evidence leak checks, tracked binary hygiene, Docker/container/port cleanup, and GitNexus detect all passed. Removed generated Python `__pycache__`. Post-slice milestone closure remains next in sequence.

## Verification

All final M039 checks passed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `json.tool + py_compile + provisioning dry-run + verify_onnx_artifacts + verify_onnx_export_contract` | 0 | ✅ pass — manifest valid and provisioning/export checks passed | 22100ms |
| 2 | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | 0 | ✅ pass — no output | 22100ms |
| 3 | `cd api && go test ./... -short` | 0 | ✅ pass — 87 passed in 4 packages | 22000ms |
| 4 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 22000ms |
| 5 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 passed in 1 package | 6300ms |
| 6 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags 'onnx hf_tokenizers' ./embed -run 'TestNewONNXEmbedderRequiresManifestPath|TestNativeHFTokenizerMatchesBaseline' -count=1` | 0 | ✅ pass — 2 passed in 1 package | 6300ms |
| 7 | `docker build -f api/Dockerfile -t fd-api:m039-default-final api` | 0 | ✅ pass — image fd-api:m039-default-final built | 6200ms |
| 8 | `M039 evidence leak checks` | 0 | ✅ pass — missing_markers=0, raw_input_or_secret_leaks=0, signed_url_like=0 | 21900ms |
| 9 | `tracked binary hygiene` | 0 | ✅ pass — tracked_native_onnx_runtime_binaries=0 | 6200ms |
| 10 | `docker ps fd-onnx-m039 and lsof port 18000` | 0 | ✅ pass — no container listed, port_18000_clean | 6100ms |
| 11 | `gitnexus_detect_changes(scope=all)` | 0 | ✅ pass — low risk, no affected processes | 0ms |
| 12 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |

## Deviations

Milestone validation/completion, DB checkpoint, commit, and GitNexus reindex run after S02 closes to preserve GSD ordering.

## Known Issues

Hosted workflow proof, exact source proof, Redis L2 restart proof, and production rollout remain open.

## Files Created/Modified

- `benchmark-results/fd-onnx-docker-smoke-m039-s01.txt`
- `benchmark-results/fd-onnx-docker-smoke-rerun-m039-s01.txt`
- `benchmark-results/fd-legal-retrieval-m039-docker-onnx-target-runtime.txt`
- `benchmark-results/fd-benchmark-m039-docker-onnx-target-runtime.txt`
- `benchmark-results/fd-onnx-docker-target-runtime-acceptance-m039-s02.txt`
