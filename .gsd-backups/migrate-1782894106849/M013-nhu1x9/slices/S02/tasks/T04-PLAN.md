---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T04: Verify build tag isolation

Run S02 verification: default tests/lint without native flags, tagged tests with native flags if implemented, artifact/leak checks, GitNexus detect_changes.

## Inputs

- `api/embed/hf_tokenizer_native.go`
- `benchmark-results/fd-tokenizer-tagged-native-m013-s02.txt`

## Expected Output

- `Task summary with verification evidence`

## Verification

Default build and applicable tagged build gates pass; no raw text/native binaries tracked.

## Observability Impact

Confirms build isolation before runtime integration work.
