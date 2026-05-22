# Requirements

This file is the explicit capability and coverage contract for the project.

## Active

### R001 — Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Class: quality-attribute
- Status: active
- Description: Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Why it matters: Latency gains are not useful if Russian legal-domain semantic quality regresses.
- Source: user clarification during M008 optimization research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S03, M040-pbp9z1/S04
- Validation: Mapped to M040 legal no-regression gate and final runtime recommendation evidence.
- Notes: M040 maps this requirement to same-host service readiness. `deepvk/USER-bge-m3` remains baseline; any runtime or candidate model must preserve Russian/legal-domain quality.

### R002 — Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Class: quality-attribute
- Status: active
- Description: Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Why it matters: During research, chunks and vectors may be reused several times; short cache retention increases model load and slows experimentation.
- Source: user clarification during M008 Redis optimization research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S04
- Validation: Mapped to M040 Docker restart harness and benchmark Redis L2 restart proof.
- Notes: M040 will prove Redis L2 reuse across packaged service restart instead of leaving restart proof skipped.

### R003 — Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Class: operability
- Status: active
- Description: Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Why it matters: Research and VPS deployment need fine tuning without rebuilding code or editing source files.
- Source: user clarification during M008 Redis/cache architecture research
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S02, M040-pbp9z1/S04
- Validation: Mapped to M040 service contract, restart harness, and final operating contract.
- Notes: M040 narrows env/config work to same-host service runtime/cache contract and required ONNX runtime verification env such as ONNX_RUNTIME_SHA256.

### R004 — Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Class: operability
- Status: active
- Description: Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Why it matters: Performance results are not comparable if env tuning differs invisibly between runs.
- Source: user clarification during M008 benchmark/config research
- Primary owning slice: M040-pbp9z1/S02
- Supporting slices: M040-pbp9z1/S04
- Validation: Mapped to M040 benchmark artifacts and final recommendation artifact.
- Notes: M040 benchmark and recommendation artifacts must retain sanitized effective config and exclude raw legal/probe text, secrets, and signed URLs.

### R005 — fd must provide a same-host local HTTP embedding service contract for neighboring services, centered on `/v1/embeddings`, batch embeddings, and `/health`.
- Class: core-capability
- Status: active
- Description: fd must provide a same-host local HTTP embedding service contract for neighboring services, centered on `/v1/embeddings`, batch embeddings, and `/health`.
- Why it matters: Neighboring services need a clear, stable local integration surface rather than an open-ended runtime experiment.
- Source: user
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S04
- Validation: Mapped to M040 service contract artifact and final recommendation.
- Notes: M040 scope excludes embedded/library API and remote hosted CI proof.

### R009 — The local embedding service must avoid silent per-request fallback between TEI and ONNX runtimes or between different tokenizers/models within one service run.
- Class: operability
- Status: active
- Description: The local embedding service must avoid silent per-request fallback between TEI and ONNX runtimes or between different tokenizers/models within one service run.
- Why it matters: Neighboring services must know which embedding contract is serving requests to avoid mixed-vector correctness issues.
- Source: inferred
- Primary owning slice: M040-pbp9z1/S01
- Supporting slices: M040-pbp9z1/S04
- Validation: Mapped to service contract and error handling strategy.
- Notes: Fallback/rollback is an operational restart/reconfiguration procedure, not hidden request-level behavior.

## Validated

### R006 — The TEI-vs-ONNX runtime recommendation must be based on an evidence envelope covering legal quality, same-host performance, restart/cache behavior, health/preflight clarity, and operational simplicity.
- Class: quality-attribute
- Status: validated
- Description: The TEI-vs-ONNX runtime recommendation must be based on an evidence envelope covering legal quality, same-host performance, restart/cache behavior, health/preflight clarity, and operational simplicity.
- Why it matters: The project goal is the best local embedding service for quality and speed, not ONNX experimentation for its own sake.
- Source: user
- Primary owning slice: M040-pbp9z1/S04
- Supporting slices: M040-pbp9z1/S01, M040-pbp9z1/S02, M040-pbp9z1/S03
- Validation: M040 S04 final artifact `benchmark-results/fd-runtime-recommendation-m040-s04.md` passed `tools/verify_m040_s04_recommendation.py` against S02/S03 evidence inputs in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, proving the final TEI-vs-ONNX evidence envelope and recommendation caveats are machine-checkable.
- Notes: Validated by S04 closeout verification; recommendation keeps TEI default until explicit ONNX switch and defers alternative model replacement fail-closed.

### R007 — M040 must not treat hosted GitHub Actions proof, remote workflow dispatch, push, upload, or artifact mirroring as required readiness gates.
- Class: constraint
- Status: validated
- Description: M040 must not treat hosted GitHub Actions proof, remote workflow dispatch, push, upload, or artifact mirroring as required readiness gates.
- Why it matters: The target deployment is same-host local service readiness, so hosted CI proof is outside the relevant acceptance boundary.
- Source: user
- Primary owning slice: M040-pbp9z1/S04
- Supporting slices: none
- Validation: M040 S04 verifier and final artifact passed in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, including hosted/remote CI readiness-gate rejection semantics and final artifact language that keeps hosted CI proof out of the same-host readiness gate.
- Notes: Validated by S04 closeout verification; readiness depends on same-host contract, preflight, cache namespace isolation, and smoke `POST /v1/embeddings`, not hosted CI.

### R008 — Alternative embedding model checks must be bounded to 1-2 plausible candidates and must use legal-domain evidence before any model can challenge `deepvk/USER-bge-m3`.
- Class: constraint
- Status: validated
- Description: Alternative embedding model checks must be bounded to 1-2 plausible candidates and must use legal-domain evidence before any model can challenge `deepvk/USER-bge-m3`.
- Why it matters: The user wants excellent legal-domain quality and speed without open-ended model experimentation.
- Source: user
- Primary owning slice: M040-pbp9z1/S03
- Supporting slices: M040-pbp9z1/S04
- Validation: M040-pbp9z1/S03 produced `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`, validated by `tools/verify_legal_model_quick_gate_artifact.py --max-candidates 2` plus closeout schema checks. The artifact caps candidates to BAAI/bge-m3 and intfloat/multilingual-e5-large, records sanitized legal-corpus hashes/counts, rejects cross-model cosine parity as a replacement criterion, and defers candidate replacement fail-closed because baseline `/health` lacks runtime metadata.
- Notes: S03 validated the bounded alternative-model gate contract. Candidate replacement remains deferred until a baseline and candidate endpoints expose contract-required runtime metadata and live legal retrieval metrics can run.

## Deferred

## Out of Scope

## Traceability

| ID | Class | Status | Primary owner | Supporting | Proof |
|---|---|---|---|---|---|
| R001 | quality-attribute | active | M040-pbp9z1/S02 | M040-pbp9z1/S03, M040-pbp9z1/S04 | Mapped to M040 legal no-regression gate and final runtime recommendation evidence. |
| R002 | quality-attribute | active | M040-pbp9z1/S02 | M040-pbp9z1/S04 | Mapped to M040 Docker restart harness and benchmark Redis L2 restart proof. |
| R003 | operability | active | M040-pbp9z1/S01 | M040-pbp9z1/S02, M040-pbp9z1/S04 | Mapped to M040 service contract, restart harness, and final operating contract. |
| R004 | operability | active | M040-pbp9z1/S02 | M040-pbp9z1/S04 | Mapped to M040 benchmark artifacts and final recommendation artifact. |
| R005 | core-capability | active | M040-pbp9z1/S01 | M040-pbp9z1/S04 | Mapped to M040 service contract artifact and final recommendation. |
| R006 | quality-attribute | validated | M040-pbp9z1/S04 | M040-pbp9z1/S01, M040-pbp9z1/S02, M040-pbp9z1/S03 | M040 S04 final artifact `benchmark-results/fd-runtime-recommendation-m040-s04.md` passed `tools/verify_m040_s04_recommendation.py` against S02/S03 evidence inputs in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, proving the final TEI-vs-ONNX evidence envelope and recommendation caveats are machine-checkable. |
| R007 | constraint | validated | M040-pbp9z1/S04 | none | M040 S04 verifier and final artifact passed in gsd_exec c52073f9-7ea0-4b13-9efa-99d54193c6f0, including hosted/remote CI readiness-gate rejection semantics and final artifact language that keeps hosted CI proof out of the same-host readiness gate. |
| R008 | constraint | validated | M040-pbp9z1/S03 | M040-pbp9z1/S04 | M040-pbp9z1/S03 produced `benchmark-results/fd-legal-model-quick-gate-m040-s03.md`, validated by `tools/verify_legal_model_quick_gate_artifact.py --max-candidates 2` plus closeout schema checks. The artifact caps candidates to BAAI/bge-m3 and intfloat/multilingual-e5-large, records sanitized legal-corpus hashes/counts, rejects cross-model cosine parity as a replacement criterion, and defers candidate replacement fail-closed because baseline `/health` lacks runtime metadata. |
| R009 | operability | active | M040-pbp9z1/S01 | M040-pbp9z1/S04 | Mapped to service contract and error handling strategy. |

## Coverage Summary

- Active requirements: 6
- Mapped to slices: 6
- Validated: 3 (R006, R007, R008)
- Unmapped active requirements: 0
