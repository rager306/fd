---
id: T02
parent: S02
milestone: M013-nhu1x9
key_files:
  - api/embed/hf_tokenizer_native.go
  - api/embed/hf_tokenizer_native_test.go
  - api/go.mod
  - api/go.sum
key_decisions:
  - `github.com/daulet/tokenizers` is added to `api/go.mod`, but imports are isolated to files guarded by `//go:build hf_tokenizers`.
  - Tagged test requires `CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64'` when run from `api`.
  - S02 only proves build-tag isolation and token parity; it does not replace `ONNXEmbedder` runtime tokenization yet.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:44:43.840Z
blocker_discovered: false
---

# T02: Implemented an opt-in `hf_tokenizers` build-tag boundary and proved default builds stay clean while tagged parity test passes.

**Implemented an opt-in `hf_tokenizers` build-tag boundary and proved default builds stay clean while tagged parity test passes.**

## What Happened

Added build-tagged native tokenizer files under `api/embed`. `hf_tokenizer_native.go` is compiled only with `-tags hf_tokenizers` and imports `github.com/daulet/tokenizers`; it wraps `tokenizers.FromFile` and returns int `input_ids` and `attention_mask`. `hf_tokenizer_native_test.go` is also tagged and compares the native tokenizer output against the M012 baseline. Default `go test ./... -short` passes without native linker flags, proving default build isolation. Tagged test passes with `CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64'`.

## Verification

Default Go tests passed without native flags. Tagged tokenizer parity test passed when linking against the validated local native artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — default build does not require native tokenizer artifact | 0ms |
| 2 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -run TestNativeHFTokenizerMatchesBaseline -count=1 -v` | 0 | ✅ pass — tagged native tokenizer parity test passed | 0ms |

## Deviations

Implemented a tagged parity test directly in Go rather than only a probe command. The test includes the fixed probe texts to exercise the native tokenizer against the persisted baseline; these texts are not rendered in artifacts.

## Known Issues

The tagged test currently duplicates fixed probe texts in test code; artifacts still exclude raw probe text. Future cleanup could centralize probe loading if needed.

## Files Created/Modified

- `api/embed/hf_tokenizer_native.go`
- `api/embed/hf_tokenizer_native_test.go`
- `api/go.mod`
- `api/go.sum`
