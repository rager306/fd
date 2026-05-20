---
id: T01
parent: S02
milestone: M013-nhu1x9
key_files:
  - api/embed/onnx.go
  - api/go.mod
  - docs/onnx-artifacts/hf-tokenizers-linux-amd64.json
key_decisions:
  - Use build tag `hf_tokenizers` for any file importing `github.com/daulet/tokenizers`.
  - Create a small internal boundary around tokenization output rather than replacing `ONNXEmbedder` directly in S02.
  - Default builds must never import `daulet/tokenizers`; tagged files only compile when `-tags hf_tokenizers` is used and `CGO_LDFLAGS=-L../.gsd/runtime/tokenizers/linux-amd64` is supplied from the `api` directory.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:42:56.643Z
blocker_discovered: false
---

# T01: Designed the `hf_tokenizers` build-tag boundary to isolate native tokenizer dependencies from default TEI builds.

**Designed the `hf_tokenizers` build-tag boundary to isolate native tokenizer dependencies from default TEI builds.**

## What Happened

Designed the build-tag boundary. S02 should add tagged-only files, likely `api/embed/hf_tokenizer_native.go` and `api/embed/hf_tokenizer_native_test.go`, with `//go:build hf_tokenizers`. Those files may import `github.com/daulet/tokenizers`; no untagged API/runtime file should import it. The tagged code should expose a small helper that loads `tokenizer.json`, returns input IDs and attention masks, and can be tested against the M012 baseline. Runtime replacement of `ONNXEmbedder` remains out of S02 unless the tagged boundary proves safe.

## Verification

Read native artifact manifest, ONNX embedder code, HF binding parity artifact, and go.mod. The summary states build tag, file names, interface shape, and default-build safety rule.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `read docs/onnx-artifacts/hf-tokenizers-linux-amd64.json` | 0 | ✅ pass — build tag and CGO_LDFLAGS contract identified | 0ms |
| 2 | `read api/embed/onnx.go and api/go.mod` | 0 | ✅ pass — current untagged tokenizer path and dependency scope identified | 0ms |
| 3 | `read benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` | 0 | ✅ pass — candidate parity evidence confirmed | 0ms |

## Deviations

None.

## Known Issues

Go module dependencies may include `daulet/tokenizers` after T02; need verify default builds still pass without `CGO_LDFLAGS`, because Go may resolve but not compile tagged files by default.

## Files Created/Modified

- `api/embed/onnx.go`
- `api/go.mod`
- `docs/onnx-artifacts/hf-tokenizers-linux-amd64.json`
