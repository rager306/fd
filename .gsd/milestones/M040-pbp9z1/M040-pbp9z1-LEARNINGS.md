---
phase: milestone-closeout
phase_name: Same-host embedding service readiness
project: fd
generated: 2026-05-22T00:00:00.000Z
counts:
  decisions: 5
  lessons: 4
  patterns: 4
  surprises: 3
missing_artifacts: []
---

# M040-pbp9z1 Learnings

### Decisions

- Keep TEI as the current/default runtime until an operator explicitly opts into ONNX; recommend packaged ONNX only for same-host performance deployments that satisfy the operating contract, artifact/tokenizer/runtime preflight, isolated Redis cache namespace, and a live smoke embedding request.
  Source: S04-SUMMARY.md/Key decisions
- Preserve `deepvk/USER-bge-m3` for the runtime recommendation and defer alternative model replacement fail-closed until bounded legal-domain evidence and runtime metadata are available.
  Source: S03-SUMMARY.md/What Happened
- Treat `/v1/embeddings` request `model` as compatibility metadata, not as a runtime/model selector; response model and `/health.runtime.model` are authoritative.
  Source: S01-SUMMARY.md/Key decisions
- Treat `/health` as operational metadata rather than live inference readiness; full runtime readiness requires a smoke `POST /v1/embeddings` probe.
  Source: S01-SUMMARY.md/Known Limitations
- Keep hosted or remote CI proof out of same-host readiness gates; readiness is based on the local contract, preflight, cache namespace isolation, and smoke embedding proof.
  Source: M040-pbp9z1-VALIDATION.md/Requirement Coverage

### Lessons

- Running Go tests from the repository root fails because the Go module lives under `api`; milestone verification must use `cd api && go test ./... -short` for module-scoped checks.
  Source: S01-SUMMARY.md/Deviations
- Redis cache namespace isolation is mandatory for TEI/ONNX or model comparisons; otherwise stale vectors can contaminate comparison evidence.
  Source: S02-SUMMARY.md/Patterns established
- Legal model replacement cannot be inferred from cross-model cosine or top-1 parity; retrieval metrics and operational compatibility are required before a candidate can challenge the baseline model.
  Source: S03-SUMMARY.md/Key decisions
- The current S03 live comparison stopped correctly because the same-host baseline `/health` lacked the contract-required runtime object; fail-closed `defer_candidate` is preferable to comparing uninspectable runtimes.
  Source: S03-SUMMARY.md/Deviations

### Patterns

- Pair final recommendation documents with semantic verifiers that fail closed on missing evidence links, weakened caveats, unsafe fallback language, and redaction violations.
  Source: S04-SUMMARY.md/Patterns established
- Use artifact-semantic verification for runtime proofs by checking health/runtime metadata, benchmark semantics, legal guard result, redaction boundaries, and cleanup state instead of trusting command exit codes alone.
  Source: S02-SUMMARY.md/Patterns established
- Publish one canonical same-host service contract and link to it from README/final recommendation documents rather than duplicating contract language across artifacts.
  Source: S01-SUMMARY.md/Patterns established
- Represent unavailable runtime evidence explicitly with stop reasons and `defer_candidate`, not implicit success, so downstream recommendations cannot overclaim model readiness.
  Source: S03-SUMMARY.md/Patterns established

### Surprises

- The packaged ONNX proof required local Python dependency repair by installing Redis into `.gsd/runtime/python-packages` because the host Python environment lacked the benchmark Redis dependency.
  Source: S02-SUMMARY.md/Deviations
- The S02 proof intentionally preserves the isolated Redis container for local cache reuse while removing only the API proof container and clearing port 18000.
  Source: S02-SUMMARY.md/Known Limitations
- GitNexus MCP tools were not directly exposed during S04, so the repo-scoped CLI equivalent was used for change-scope auditing.
  Source: S04-SUMMARY.md/Deviations
