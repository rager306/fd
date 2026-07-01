---
verdict: pass
remediation_round: 0
---

# Milestone Validation: M011-33b7wf

## Success Criteria Checklist
- [x] ONNX backend remains opt-in and TEI remains default — `EMBEDDING_BACKEND` defaults to `tei`; ONNX requires explicit env vars.
- [x] Artifact checksum/config validation exists before ONNX load — manifest validator checks path, size, SHA256, output metadata, dimensions, and normalization expectations.
- [x] Dense ONNX backend prototype either works locally or is blocked with evidence — Go ONNX Runtime loads/runs, but semantic equivalence is blocked by tokenizer mismatch; evidence is persisted.
- [x] TEI vs ONNX comparator evidence is persisted — `benchmark-results/fd-go-onnx-m011-s03.txt` records failed isolated-cache comparison and no raw probe text.
- [x] Performance benchmarking is intentionally deferred — valid because semantic equivalence failed and speed claims would be misleading.
- [x] No large ONNX artifacts are committed — tracked large artifact scan passed.
- [x] No production runtime switch, model replacement, INT8, provider variants, or rewrite occurred — TEI remains default and ONNX stays opt-in.

## Slice Delivery Audit
| Slice | Claimed output | Delivered output | Verdict |
|---|---|---|---|
| S01 | Artifact manifest and checksum contract | `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` with local artifact checksum/metadata; local binary remains ignored | Pass |
| S02 | Opt-in backend seam | `EMBEDDING_BACKEND` config, manifest validation, explicit ONNX env requirements, TEI default preserved | Pass |
| S03 | ONNX backend works locally or blocks with evidence | Go ONNX runtime loads/runs; isolated-cache comparison failed; tokenizer mismatch evidence captured | Pass as blocked-with-evidence |
| S04 | Recommendation and safety verification | S04 research recommends tokenizer parity next; fresh Go/lint/Compose/health/manifest/GitNexus gates passed | Pass |

## Cross-Slice Integration
- S01 manifest contract feeds S02 config validation and S03 ONNX embedder construction.
- S02 backend seam preserves TEI default and requires explicit ONNX manifest/runtime/tokenizer env vars.
- S03 implements Go ONNX load/run and then supplies blocker evidence to S04.
- S04 consumes S03 evidence and correctly prevents invalid throughput benchmarking.
- No boundary mismatch remains: the milestone allows a prototype to either work or be blocked with evidence.

## Requirement Coverage
- Existing requirement to preserve TEI default: satisfied; TEI remains default and health check passed.
- Existing requirement to avoid model replacement: satisfied; exact `deepvk/USER-bge-m3` artifact/manifest used.
- Existing requirement to avoid committing large artifacts: satisfied; tracked large model artifact scan found zero `.onnx`/`.safetensors` files.
- New requirement surfaced: tokenizer parity must be validated before ONNX performance benchmarking or production recommendation.

## Verification Class Compliance
- Unit/integration: `cd api && go test ./... -short` passed with 78 tests across 4 packages.
- Lint: pinned GolangCI-Lint v2.12.2 reported 0 issues.
- Runtime safety: Docker Compose config rendered and default API health returned ok.
- Artifact integrity: manifest JSON parsed; local ONNX artifact size/SHA256 matched manifest.
- Evidence hygiene: comparison artifact has raw_probe_texts_logged=false; no tracked large model artifacts.
- Code intelligence: GitNexus detect_changes reported low risk and no affected processes.
- LSP diagnostics unavailable because no Go language server is configured.


## Verdict Rationale
M011 satisfies its safe-prototype contract: it either had to produce a working opt-in ONNX prototype or block with concrete evidence. It produced a loadable/runnable Go ONNX path but failed semantic equivalence due tokenizer parity; that blocker is documented and verified, while the default TEI path remains safe.
