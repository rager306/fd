# S04: Alternative model research

**Goal:** Research candidate embedding models that may perform better in this Russian/legal environment and define how to benchmark them fairly.
**Demo:** After this, alternative models are ranked for future legal-corpus benchmarking.

## Must-Haves

- Candidate models are ranked.
- Each candidate lists rationale, risks, license/deployability, and benchmark requirements.
- MiniLM-like references are not treated as replacements without legal evidence.
- Future benchmark protocol uses the 44-ФЗ gate and sanitized config snapshots.

## Proof Level

- This slice proves: Research artifact and benchmark protocol.

## Integration Closure

Creates a bounded candidate list without derailing current USER-bge-m3 root-cause work.

## Verification

- Records source links, licenses, dimensions, ONNX/CPU feasibility, expected memory, and legal/Russian evidence gaps.

## Tasks

- [x] **T01: Research candidate embedding models** `est:medium`
  Research current Russian/multilingual/legal embedding model candidates and collect source evidence for suitability, dimensions, license, and ONNX/CPU deployment feasibility.
  - Verify: At least 5 candidates with source links and caveats are identified.

- [x] **T02: Write model alternatives research artifact** `est:small`
  Write `benchmark-results/fd-model-alternatives-m016-s04.txt` with ranked candidates and a fair benchmark protocol using the 44-ФЗ gate.
  - Files: `benchmark-results/fd-model-alternatives-m016-s04.txt`
  - Verify: Artifact exists, cites sources, ranks candidates, and avoids production switch claims.

## Files Likely Touched

- benchmark-results/fd-model-alternatives-m016-s04.txt
