---
id: T04
parent: S03
milestone: M013-nhu1x9
key_files:
  - api/embed/onnx.go
  - api/embed/onnx_tokenizer_default.go
  - api/embed/onnx_tokenizer_hf.go
  - api/embed/hf_tokenizer_native.go
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
key_decisions:
  - S03 integration is accepted: tagged ONNX cosine passes, so ONNX performance benchmarking is now meaningful for the tagged HF tokenizer path only.
  - Default build remains green and untagged runtime behavior is preserved.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:55:40.266Z
blocker_discovered: false
---

# T04: Verified tagged ONNX integration: default build safe, tagged tests pass, cosine passes, no leaks or lingering servers.

**Verified tagged ONNX integration: default build safe, tagged tests pass, cosine passes, no leaks or lingering servers.**

## What Happened

Ran final S03 verification. Default tests and lint pass. Tagged embed tests pass with the native tokenizer library. The cosine artifact parses, records the isolated cache namespace, passes, and contains no raw probe text. No native binaries are tracked and no tagged background server remains. GitNexus reports medium risk limited to ONNX embedder internal processes, which are covered by tagged tests and the successful cosine comparison.

## Verification

Fresh verification passed: default tests/lint, tagged tests, artifact hygiene, no background processes, and GitNexus scope reviewed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — 78 passed in 4 packages | 10400ms |
| 2 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 10300ms |
| 3 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -count=1` | 0 | ✅ pass — 20 tagged embed tests passed | 10300ms |
| 4 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python artifact hygiene check` | 0 | ✅ pass — s03_cosine_artifact_check=pass; raw_probe_text_leaks=0; tracked_native_binaries=0 | 0ms |
| 5 | `bg_shell list` | 0 | ✅ pass — no background processes | 0ms |
| 6 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ⚠️ medium — ONNX embedder internal processes affected and verified | 0ms |

## Deviations

GitNexus reports medium scope for ONNX embedder internals, which matches the intentional tagged tokenizer integration. Fresh default and tagged tests plus cosine comparison cover the affected path.

## Known Issues

Tagged native library packaging is still local; Docker/CI tagged builds are not yet added. Larger Russian/legal corpus validation remains future work before production switch.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/embed/onnx_tokenizer_default.go`
- `api/embed/onnx_tokenizer_hf.go`
- `api/embed/hf_tokenizer_native.go`
- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
