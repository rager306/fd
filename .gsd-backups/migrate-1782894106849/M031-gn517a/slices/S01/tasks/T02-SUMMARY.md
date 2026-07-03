---
id: T02
parent: S01
milestone: M031-gn517a
key_files:
  - .gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md
  - .gsd/exec/78224446-6bc5-4bcb-84f1-633ade016a0c.stdout
  - .gsd/exec/24907b45-1ff7-436e-a6fa-b99bc5be5f06.stdout
key_decisions:
  - Pinned `daulet/tokenizers` `v1.27.0` release asset as the native tokenizer source candidate because the extracted library matches the existing local checksum.
  - Pinned HF revision URL for tokenizer.json as the tokenizer JSON source candidate because it matches manifest size/sha.
  - Pinned PyPI onnxruntime 1.26.0 CP313 manylinux wheel as ONNX Runtime source candidate because the extracted library matches the existing local checksum.
  - Kept exported ONNX model binary blocked until it is mirrored/uploaded to an immutable non-secret source or a separate reproducible-export path is validated.
duration: 
verification_result: passed
completed_at: 2026-05-21T06:33:38.454Z
blocker_discovered: false
---

# T02: Assessed immutable source candidates and kept the ONNX model binary blocker explicit.

**Assessed immutable source candidates and kept the ONNX model binary blocker explicit.**

## What Happened

Queried public metadata for Hugging Face, GitHub releases, and PyPI. Verified that `daulet/tokenizers` v1.27.0 contains the exact local `libtokenizers.a` checksum, the pinned Hugging Face revision URL contains the expected tokenizer.json checksum, and the PyPI onnxruntime 1.26.0 CP313 manylinux wheel contains the expected `libonnxruntime.so.1.26.0` checksum. The exported ONNX model remains blocked because no immutable source for the exact binary exists.

## Verification

Research notes distinguish immutable source candidates from blockers and cite public metadata/checksum probe artifacts.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gsd_exec public metadata probe for PyPI, GitHub releases, HF pinned URLs` | 0 | ✅ pass — found candidate source metadata | 12843ms |
| 2 | `gsd_exec native tokenizer v1.27.0 archive checksum probe` | 0 | ✅ pass — extracted libtokenizers.a matches expected sha256 | 895ms |
| 3 | `gsd_exec HF tokenizer JSON and ONNX Runtime wheel checksum probe` | 0 | ✅ pass — tokenizer.json and libonnxruntime.so.1.26.0 match expected sha256 | 62056ms |

## Deviations

Performed bounded downloads of small public artifacts to match exact checksums for native tokenizer, tokenizer JSON, and ONNX Runtime. This did not download the large ONNX model binary or change external state.

## Known Issues

ONNX model binary remains local-only and blocked for hosted proof. Source candidate acceptance still needs docs/manifest wiring in S02 and a later real hosted workflow proof after push approval.

## Files Created/Modified

- `.gsd/milestones/M031-gn517a/slices/S01/S01-RESEARCH.md`
- `.gsd/exec/78224446-6bc5-4bcb-84f1-633ade016a0c.stdout`
- `.gsd/exec/24907b45-1ff7-436e-a6fa-b99bc5be5f06.stdout`
