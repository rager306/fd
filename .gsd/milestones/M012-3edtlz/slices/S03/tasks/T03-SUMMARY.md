---
id: T03
parent: S03
milestone: M012-3edtlz
key_files:
  - benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt
  - tools/compare_tokenizers.py
key_decisions:
  - Do not import `daulet/tokenizers` into default API code in this slice.
  - Treat native library packaging/build tags as the next required design gate before runtime integration.
  - S04 should recommend a build-tag or packaging milestone if continuing pure Go ONNX integration.
duration: 
verification_result: passed
completed_at: 2026-05-20T02:20:07.478Z
blocker_discovered: true
---

# T03: Recorded the runtime integration blocker: HF Rust tokenizer binding passes parity, but native library packaging/build tags are required before API integration.

**Recorded the runtime integration blocker: HF Rust tokenizer binding passes parity, but native library packaging/build tags are required before API integration.**

## What Happened

Evaluated whether to integrate the parity-passing HF Rust tokenizers binding into `api/embed/onnx.go`. Although token-level parity passed, integration would require adding `github.com/daulet/tokenizers` and linking `libtokenizers.a`. Because Go linking happens at build time, importing this dependency into normal API code would make default builds depend on a native artifact that is not yet packaged in Docker or CI. The safe outcome for S03 is therefore to keep runtime code unchanged and record the integration blocker. The candidate is correct, but packaging/build-tag design is required before integration.

## Verification

No runtime code was changed. The decision is backed by the passing candidate artifact and known native `libtokenizers.a` link requirement from the isolated probe.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt` | 0 | ✅ parity evidence — all five probes pass HF baseline comparison | 0ms |
| 2 | `isolated `daulet/tokenizers` probe required CGO_LDFLAGS=-L/tmp/fd-daulet-tokenizers-probe and libtokenizers.a` | 0 | ⚠️ packaging blocker — native static library required at build/link time | 0ms |
| 3 | `runtime integration skipped` | 0 | ✅ safe — default API build not changed by native dependency | 0ms |

## Deviations

Parity passed, but runtime integration was intentionally not performed because importing `github.com/daulet/tokenizers` into default API code would introduce a native static-library link requirement and likely break normal builds/CI unless packaging and build tags are designed first.

## Known Issues

The only parity-passing Go path currently requires `libtokenizers.a` and CGO/linker configuration. Project Dockerfiles/CI do not yet provide this native artifact, so direct integration would be unsafe.

## Files Created/Modified

- `benchmark-results/fd-tokenizer-go-hf-binding-m012-s03.txt`
- `tools/compare_tokenizers.py`
