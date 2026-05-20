# S02: Opt in build tag boundary — UAT

**Milestone:** M013-nhu1x9
**Written:** 2026-05-20T03:47:19.200Z

# S02 UAT — Opt in build tag boundary

## Checks

- [x] Native tokenizer imports are guarded by `//go:build hf_tokenizers`.
- [x] Default `go test ./... -short` passes without native flags.
- [x] Default pinned lint passes without native flags.
- [x] Tagged native tokenizer parity test passes with `CGO_LDFLAGS`.
- [x] Tagged parity artifact exists and records PASS.
- [x] No native binary is tracked.
- [x] Artifact leak checks pass.
- [x] GitNexus reports low risk and no affected processes.

## UAT Result

Pass. S03 may integrate the tagged tokenizer path into the opt-in ONNX backend.

