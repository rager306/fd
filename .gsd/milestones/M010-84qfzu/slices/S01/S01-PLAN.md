# S01: ONNX artifact provenance

**Goal:** Find the safest BGE-M3 FP32 dense-only ONNX artifact or export path and record provenance requirements.
**Demo:** After this, candidate ONNX/export paths are ranked by provenance, artifact availability, and dense-output compatibility risk.

## Must-Haves

- Candidate ONNX/export paths are listed and ranked.
- Required artifacts are identified: model ONNX, tokenizer, config, output names, hashes.
- Dense-only compatibility risks are documented.
- No model replacement is recommended.

## Proof Level

- This slice proves: source/web/model-card evidence plus local filesystem checks

## Integration Closure

Feeds comparator/export decisions for S02/S03.

## Verification

- Defines artifact hash/config fields that future benchmark snapshots must record.

## Tasks

- [x] **T01: Inspect local model artifact storage** `est:small`
  Inspect local TEI/Docker model storage and project config to determine what model artifacts are already present locally and where large ONNX artifacts would live without being committed.
  - Files: `docker-compose.yaml`, `docker-compose.override.yaml`, `.gitignore`
  - Verify: Local Docker volumes/config inspected; no large artifacts staged.

- [x] **T02: Research ONNX artifact candidates** `est:medium`
  Research BGE-M3 ONNX artifact candidates and export paths from current sources: official model/export docs, community ONNX model cards, TEI/optimum/export options, and Go/Rust wrapper relevance. Rank candidates by provenance and dense-output compatibility.
  - Verify: Research cites source URLs and ranks candidates.

- [x] **T03: Synthesize ONNX provenance path** `est:small`
  Save S01 research synthesis with candidate ranking, required artifacts, dense-output risks, hash/provenance fields, and go/no-go input for S02/S03.
  - Verify: S01 research artifact exists and states no production runtime change.

## Files Likely Touched

- docker-compose.yaml
- docker-compose.override.yaml
- .gitignore
