---
id: T03
parent: S01
milestone: M017-j10hmp
key_files:
  - benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt
key_decisions:
  - Tagged Go ONNX 512 improves retrieval parity dramatically versus M015 128, but does not pass strict cross-backend cosine threshold due very long fragments above 512 tokens.
  - S02 must recommend chunking or longer-sequence handling rather than treating 512 alone as sufficient.
duration: 
verification_result: mixed
completed_at: 2026-05-20T07:27:46.868Z
blocker_discovered: false
---

# T03: Ran the full ONNX 512 legal quality gate and captured a measured strict FAIL with strong ranking parity.

**Ran the full ONNX 512 legal quality gate and captured a measured strict FAIL with strong ranking parity.**

## What Happened

Ran the full legal retrieval evaluator against TEI and tagged Go ONNX configured with max sequence length 512. The run produced a sanitized artifact and a measured FAIL verdict under the strict cosine threshold. Compared with M015, ranking parity improved strongly: top-1 agreement reached 1.0, mean overlap@5 reached 0.997701, and ONNX recall ratio reached 1.0. However, the minimum cross-backend cosine remained 0.98982302 because the longest legal fragments still exceed 512 tokens.

## Verification

Evaluator wrote the artifact, exited 2 due measured FAIL verdict, and artifact hygiene check passed with no raw legal text leaks.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `uv run --python 3.13 --with requests python tools/evaluate_legal_retrieval.py ... --onnx-runtime-label tagged-onnx-hf-512 --onnx-cache-namespace m017-onnx-512-legal-quality` | 2 | ⚠️ measured FAIL — strict cosine threshold not met; artifact written | 54200ms |
| 2 | `python artifact hygiene check for benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt` | 0 | ✅ pass — m017_s01_artifact_hygiene=pass; raw_legal_text_leaks=0 | 0ms |

## Deviations

Evaluator exited with code 2 because its strict quality verdict was FAIL. This is an expected measured gate outcome, not a command/runtime blocker.

## Known Issues

Strict gate failed: minimum cross-backend cosine is 0.98982302, below threshold 0.999. Ranking metrics are strong: top1_agreement 1.0, mean_overlap_at_5 0.997701, recall ratio 1.0.

## Files Created/Modified

- `benchmark-results/fd-legal-retrieval-m017-s01-onnx512.txt`
