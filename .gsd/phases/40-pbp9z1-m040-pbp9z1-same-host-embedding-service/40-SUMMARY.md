---
id: M040-pbp9z1
title: "Same-host embedding service readiness"
status: complete
completed_at: 2026-05-22T08:40:23.433Z
key_decisions:
  - TEI remains the current/default runtime until an operator explicitly opts into packaged ONNX.
  - Packaged ONNX is recommended only for same-host performance deployments satisfying the operating contract, preflight, cache namespace isolation, and smoke embedding checks.
  - deepvk/USER-bge-m3 remains the model baseline; alternative model replacement is deferred fail-closed.
  - /v1/embeddings request model is compatibility metadata, not a per-request model/runtime selector.
  - /health exposes safe operational metadata but is not a live inference readiness probe.
  - Hosted/remote CI proof is not a readiness gate for same-host local service deployment.
key_files:
  - docs/same-host-embedding-service-contract.md
  - README.md
  - api/handlers/health.go
  - api/handlers/health_test.go
  - api/main.go
  - api/handlers/embeddings_integration_test.go
  - tools/run_m040_s02_docker_restart_proof.sh
  - tools/verify_m040_s02_artifacts.py
  - tools/evaluate_legal_model_quick_gate.py
  - tools/verify_legal_model_quick_gate_artifact.py
  - tools/verify_m040_s04_recommendation.py
  - benchmark-results/fd-runtime-recommendation-m040-s04.md
  - benchmark-results/fd-benchmark-m040-s02-onnx-docker-restart.txt
  - benchmark-results/fd-legal-model-quick-gate-m040-s03.md
  - .gsd/milestones/M040-pbp9z1/M040-pbp9z1-LEARNINGS.md
lessons_learned:
  - Run Go tests from the api module rather than repository root for fd-api verification.
  - Cache namespace isolation is mandatory for TEI/ONNX or model comparisons to avoid stale vector contamination.
  - Alternative embedding-model replacement needs legal-domain retrieval evidence and operational compatibility, not cross-model cosine parity.
  - Semantic artifact verifiers are useful closeout guards for recommendation docs because they catch missing evidence, weakened caveats, unsafe fallback language, and redaction regressions.
  - Fail-closed defer_candidate artifacts provide a safe way to record unavailable runtime evidence without overclaiming readiness.
---

# M040-pbp9z1: Same-host embedding service readiness

**Prepared fd as a same-host local HTTP embedding service by publishing the service contract, proving packaged ONNX restart/cache behavior, bounding alternative-model exploration, and issuing a machine-verified TEI-vs-ONNX recommendation.**

## What Happened

M040 reframed the embedding runtime work around same-host service readiness rather than ONNX experimentation for its own sake. S01 established the canonical local HTTP service contract for neighboring services, linked it from README, added safe runtime metadata to /health, documented /v1/embeddings model-field compatibility semantics, and made clear that /health is metadata rather than a live inference probe. S02 turned the packaged ONNX restart/cache proof into reproducible artifacts: a proof runner, verifier, preflight output, restart/cache benchmark, legal retrieval guard, and cleanup audit showing API-only restart with Redis L2 preserved. S03 bounded alternative-model exploration to exactly two candidates, rejected cross-model cosine parity as a replacement basis, and produced a sanitized fail-closed defer_candidate artifact because live runtime/legal evidence was unavailable. S04 consumed the S01/S02/S03 evidence envelope to publish benchmark-results/fd-runtime-recommendation-m040-s04.md and a semantic verifier. The final stance keeps TEI as the current/default runtime, allows explicit opt-in packaged ONNX for same-host performance deployments that satisfy preflight/cache/smoke checks, preserves deepvk/USER-bge-m3, excludes hosted/remote CI proof from same-host readiness gates, and defers alternative model replacement until bounded legal-domain evidence exists.

## Success Criteria Results

- [x] Same-host local HTTP embedding service contract exists and is grounded in current code/runtime evidence: S01 produced docs/same-host-embedding-service-contract.md, README discoverability, safe /health.runtime metadata, and /v1/embeddings compatibility semantics with Go test evidence.
- [x] Packaged Docker restart/cache behavior is proven: S02 produced tools/run_m040_s02_docker_restart_proof.sh, tools/verify_m040_s02_artifacts.py, preflight/benchmark/legal/audit artifacts, and verifier PASS for API-only restart with Redis L2 reuse and cleanup.
- [x] Legal-domain quality remains no-regression for included recommendations: S02 legal retrieval guard passed for packaged ONNX with deepvk/USER-bge-m3; S03/S04 defer unproven alternative models fail-closed.
- [x] Bounded alternative model quick gate completed: S03 evaluated exactly BAAI/bge-m3 and intfloat/multilingual-e5-large, enforced --max-candidates 2, sanitized corpus metadata, and ended with defer_candidate.
- [x] Final TEI-vs-ONNX recommendation produced from the evidence envelope: S04 published benchmark-results/fd-runtime-recommendation-m040-s04.md and verified it with tools/verify_m040_s04_recommendation.py.

## Definition of Done Results

- Duplicate completion guard passed: gsd_milestone_status showed M040-pbp9z1 active with all four slices complete before closeout.
- Code-change verification passed: merge-base diff to master had no current non-.gsd diff, so milestone-scoped commit evidence was inspected and showed non-.gsd implementation artifacts including api/handlers/health.go, api/main.go, api/main_test.go, docs/same-host-embedding-service-contract.md, tools/run_m040_s02_docker_restart_proof.sh, and tools/verify_m040_s02_artifacts.py.
- Slice completion passed: S01, S02, S03, and S04 are [x] in the roadmap and gsd_milestone_status reports all tasks complete.
- Summary artifacts passed: all four slice SUMMARY.md files and the validation artifact exist.
- Integration/semantic verification passed in gsd_exec 0cda6481-e8bc-479e-a85f-227555cce8f2: S04 recommendation verifier PASS, S02 artifact verifier PASS, S03 quick-gate artifact valid, and GitNexus detect-changes reported no changes detected.
- Roadmap checklist passed: no unchecked roadmap items and no Horizontal Checklist was present.

## Requirement Outcomes

- R001 covered: deepvk/USER-bge-m3 is preserved; S02 legal guard passed for packaged ONNX and S03/S04 reject unproven model replacement.
- R002 covered: S02 proved Redis L2 reuse across API-only restart with an isolated cache namespace; S04 requires namespace isolation.
- R003 covered: S01/S04 document runtime/env/cache knobs including explicit ONNX opt-in, preflight, ONNX_RUNTIME_SHA256, and EMBEDDING_CACHE_VERSION.
- R004 covered: S02/S04 artifacts record sanitized effective config and S04 verifier checks redaction boundaries.
- R005 covered: S01 produced the same-host local HTTP service contract for /health, /v1/embeddings, and /embeddings/batch.
- R006 validated: S04 recommendation artifact passed semantic verification against the full evidence envelope.
- R007 validated: S04 final artifact/verifier reject hosted/remote CI proof as a same-host readiness gate.
- R008 validated: S03 bounded candidate checks to two candidates with sanitized legal-domain evidence and defer_candidate.
- R009 covered: S01/S04 establish no request-level runtime/model fallback and require runtime identity plus smoke readiness.

## Deviations

S01's first root-level Go test invocation failed because the Go module is under api; module-scoped verification passed. S02 required local Python Redis dependency repair under .gsd/runtime/python-packages before the proof could run. S03 did not run live candidate metrics because baseline /health lacked runtime metadata and therefore correctly deferred fail-closed. S04 used the repo-scoped GitNexus CLI equivalent because MCP GitNexus tools were not directly exposed in that harness.

## Follow-ups

Future work may rerun the bounded S03 legal model gate after baseline and candidate endpoints expose contract-required /health.runtime metadata and live legal retrieval metrics. Operators enabling ONNX should perform artifact/tokenizer/runtime preflight, configure cache namespace isolation, set ONNX_RUNTIME_SHA256 when enforcing runtime-library integrity, and run a live /v1/embeddings smoke check.
