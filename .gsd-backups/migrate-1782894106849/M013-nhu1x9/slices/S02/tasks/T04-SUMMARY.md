---
id: T04
parent: S02
milestone: M013-nhu1x9
key_files:
  - api/embed/hf_tokenizer_native.go
  - api/embed/hf_tokenizer_native_test.go
  - benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt
  - api/go.mod
  - api/go.sum
key_decisions:
  - Default build isolation is proven: tests and lint pass without `CGO_LDFLAGS` or build tags.
  - Tagged native tokenizer path is proven: tagged test passes when `CGO_LDFLAGS` points to the validated native artifact.
  - No native binaries are tracked.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:46:45.983Z
blocker_discovered: false
---

# T04: Verified build-tag isolation: default builds stay clean and the opt-in native tokenizer parity test passes with the validated library.

**Verified build-tag isolation: default builds stay clean and the opt-in native tokenizer parity test passes with the validated library.**

## What Happened

Ran final S02 verification. Default `go test ./... -short` and pinned lint pass without native library flags. The tagged native tokenizer test passes with `CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64'` and `-tags hf_tokenizers`. The tagged artifact parses, records `raw_probe_texts_logged=false`, has no raw probe text leaks, and no `.a/.so/.dylib` files are tracked. GitNexus reports low risk with no affected processes.

## Verification

Fresh verification passed: default tests/lint, tagged test, artifact hygiene, tracked-native scan, and GitNexus detect_changes.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && go test ./... -short` | 0 | ✅ pass — 78 passed in 4 packages | 6700ms |
| 2 | `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...` | 0 | ✅ pass — 0 issues | 6600ms |
| 3 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -run TestNativeHFTokenizerMatchesBaseline -count=1` | 0 | ✅ pass — 1 tagged test passed | 6600ms |
| 4 | `uv run --python 3.13 --with transformers --with torch --with sentencepiece python artifact hygiene check` | 0 | ✅ pass — s02_tagged_artifact_check=pass; raw_probe_text_leaks=0; tracked_native_binaries=0 | 0ms |
| 5 | `gitnexus_detect_changes(scope=all, repo=fd)` | 0 | ✅ pass — low risk; affected_processes=[] | 0ms |

## Deviations

The first artifact hygiene check was accidentally run outside the `uv` dependency context and failed to import `transformers`; it was rerun under the same `uv` dependency context as tokenizer tooling and passed.

## Known Issues

Tagged test duplicates fixed probe texts in code; artifacts do not leak raw texts. Runtime `ONNXEmbedder` still uses `sugarme` until S03 integration changes it.

## Files Created/Modified

- `api/embed/hf_tokenizer_native.go`
- `api/embed/hf_tokenizer_native_test.go`
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
- `api/go.mod`
- `api/go.sum`
