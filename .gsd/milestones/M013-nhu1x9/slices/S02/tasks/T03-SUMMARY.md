---
id: T03
parent: S02
milestone: M013-nhu1x9
key_files:
  - benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt
  - api/embed/hf_tokenizer_native_test.go
key_decisions:
  - `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt` is the S02 proof artifact for project-local tagged native parity.
  - The artifact records command, manifest, native artifact checksum, PASS output, and `raw_probe_texts_logged=false`.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:45:25.056Z
blocker_discovered: false
---

# T03: Persisted the tagged native tokenizer parity proof artifact; the project-local `hf_tokenizers` path passes.

**Persisted the tagged native tokenizer parity proof artifact; the project-local `hf_tokenizers` path passes.**

## What Happened

Ran the project-local tagged native tokenizer parity test with the validated local `libtokenizers.a`, then wrote `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`. The artifact records the exact tagged command, native artifact manifest and SHA256, baseline path, raw-text exclusion flag, test output, and PASS verdict. This proves the tagged package path can reproduce HF tokenization in the project, not only in the earlier temp module.

## Verification

Tagged Go test exited 0 and artifact was written with PASS verdict and raw-text exclusion flag.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd api && CGO_LDFLAGS='-L../.gsd/runtime/tokenizers/linux-amd64' go test -tags hf_tokenizers ./embed -run TestNativeHFTokenizerMatchesBaseline -count=1 -v` | 0 | ✅ pass — tagged native tokenizer parity test passed | 0ms |
| 2 | `write benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt` | 0 | ✅ pass — artifact written with PASS verdict | 0ms |

## Deviations

Used the tagged Go test as the parity source and rendered its output into an artifact rather than extending the Python comparator again. This proves the project-local tagged package path, which is the S02 objective.

## Known Issues

Tagged parity artifact proves tokenization only; ONNX runtime still has not been rewired to use the tagged tokenizer path.

## Files Created/Modified

- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`
- `api/embed/hf_tokenizer_native_test.go`
