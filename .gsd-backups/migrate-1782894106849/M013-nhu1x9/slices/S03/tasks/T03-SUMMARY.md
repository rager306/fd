---
id: T03
parent: S03
milestone: M013-nhu1x9
key_files:
  - benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt
  - api/embed/onnx.go
  - api/embed/onnx_tokenizer_hf.go
key_decisions:
  - Tagged ONNX runtime with HF native tokenizer is semantically equivalent on fixed probes at the established `0.999` cosine threshold.
  - Use `EMBEDDING_CACHE_VERSION=m013-hf-tokenizer` for the comparison to avoid TEI cache masking.
  - Performance benchmarking is now meaningful only for this tagged/native ONNX path, not the untagged `sugarme` path.
duration: 
verification_result: passed
completed_at: 2026-05-20T03:54:34.973Z
blocker_discovered: false
---

# T03: Ran the tagged ONNX API with the HF native tokenizer and passed TEI-vs-ONNX cosine equivalence on all fixed probes.

**Ran the tagged ONNX API with the HF native tokenizer and passed TEI-vs-ONNX cosine equivalence on all fixed probes.**

## What Happened

Started the fd API locally on port 18000 with `go run -tags hf_tokenizers`, validated native `libtokenizers.a` linker flags, opt-in ONNX backend config, and isolated Redis namespace `m013-hf-tokenizer`. The tagged ONNX API started successfully. Comparing the default TEI API on port 8000 against the tagged ONNX API on port 18000 produced cosine values around `0.999993` for all five fixed probes, all above the `0.999` threshold. The comparison artifact was saved and the local tagged server was stopped.

## Verification

Tagged ONNX server started, isolated-cache comparison passed, artifact written with PASS, and server cleanup completed.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `bg_shell start tagged ONNX API on port 18000 with CGO_LDFLAGS and EMBEDDING_CACHE_VERSION=m013-hf-tokenizer` | 0 | ✅ pass — server ready | 0ms |
| 2 | `python3 TEI default vs tagged Go ONNX comparison writing benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt` | 0 | ✅ pass — all five cosines >= 0.999; observed ~0.999993 | 0ms |
| 3 | `bg_shell kill tagged ONNX API server` | 0 | ✅ cleanup | 0ms |

## Deviations

None. The tagged runtime started successfully and cosine comparison passed, so S03 did not need to record a runtime blocker.

## Known Issues

This is fixed-probe cosine only. Larger Russian/legal corpus quality and production native artifact packaging/CI still remain future gates before production switch.

## Files Created/Modified

- `benchmark-results/fd-go-onnx-hf-tokenizer-m013-s03.txt`
- `api/embed/onnx.go`
- `api/embed/onnx_tokenizer_hf.go`
