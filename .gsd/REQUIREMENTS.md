# Requirements

This file is the explicit capability and coverage contract for the project.

## Active

### R001 — Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Class: quality-attribute
- Status: active
- Description: Embedding runtime optimizations must preserve Russian-language and legal-domain retrieval/embedding quality for the current model; any model replacement requires benchmark evidence on a Russian legal corpus.
- Why it matters: Latency gains are not useful if Russian legal-domain semantic quality regresses.
- Source: user clarification during M008 optimization research
- Primary owning slice: M008-6hnowu
- Validation: A future optimization spike must either keep the current model or include a Russian legal corpus benchmark covering quality plus latency before changing the model default.
- Notes: MiniLM/all-MiniLM and other analog models may be used as implementation references only. They are not acceptable replacements unless evaluated against a Russian legal corpus with agreed quality metrics and compared to the current BGE-M3-style runtime.

### R002 — Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Class: quality-attribute
- Status: active
- Description: Research/chunking workflows must use a sufficiently long-lived embedding cache so repeated chunk processing can reuse vectors and reduce model load.
- Why it matters: During research, chunks and vectors may be reused several times; short cache retention increases model load and slows experimentation.
- Source: user clarification during M008 Redis optimization research
- Primary owning slice: M008-6hnowu
- Validation: A future cache optimization spike must expose/define embedding cache retention policy and benchmark repeated chunk reuse without unnecessary model calls.
- Notes: Current Redis TTL is 24h in api/cache/redis.go. M008 Redis research should evaluate longer configurable TTLs, no-expire research mode, model/version-aware cache keys, explicit invalidation, and memory/eviction policies. Cache duration should account for repeated legal corpus chunking experiments where the same chunks may be embedded multiple times over days or weeks.

### R003 — Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Class: operability
- Status: active
- Description: Performance/cache/runtime tuning parameters should be configurable through environment variables with safe defaults and validation.
- Why it matters: Research and VPS deployment need fine tuning without rebuilding code or editing source files.
- Source: user clarification during M008 Redis/cache architecture research
- Primary owning slice: M008-6hnowu
- Validation: A future implementation spike must document env vars, defaults, validation behavior, and which vars affect cache key namespace versus runtime-only tuning.
- Notes: Candidate env knobs include embedding cache TTL, cache namespace/version, model revision/key salt, Redis pool size/min idle/timeouts, local L1 size/TTL, batch cache mode, and benchmark/research mode. Values that affect cache correctness must be reflected in cache keys or require explicit invalidation.

### R004 — Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Class: operability
- Status: active
- Description: Benchmark artifacts must record the effective environment/configuration parameters used for the run so results remain comparable across tuning experiments.
- Why it matters: Performance results are not comparable if env tuning differs invisibly between runs.
- Source: user clarification during M008 benchmark/config research
- Primary owning slice: M008-6hnowu
- Validation: Future benchmark changes must emit or save a sanitized effective configuration section before numeric results, and benchmark comparisons must cite matching or intentionally varied config fields.
- Notes: Benchmark output should include a sanitized env/config snapshot: model id/revision, cache version/namespace, TTL/no-expire mode, Redis address/pool/timeouts, L1 cache size/TTL, batch cache mode, benchmark corpus/chunk settings, Docker image/config identifiers where available, and git commit. Secrets must never be printed. Configuration values that affect cache correctness should be clearly separated from runtime-only tuning knobs.

## Validated

## Deferred

## Out of Scope

## Traceability

| ID | Class | Status | Primary owner | Supporting | Proof |
|---|---|---|---|---|---|
| R001 | quality-attribute | active | M008-6hnowu | none | A future optimization spike must either keep the current model or include a Russian legal corpus benchmark covering quality plus latency before changing the model default. |
| R002 | quality-attribute | active | M008-6hnowu | none | A future cache optimization spike must expose/define embedding cache retention policy and benchmark repeated chunk reuse without unnecessary model calls. |
| R003 | operability | active | M008-6hnowu | none | A future implementation spike must document env vars, defaults, validation behavior, and which vars affect cache key namespace versus runtime-only tuning. |
| R004 | operability | active | M008-6hnowu | none | Future benchmark changes must emit or save a sanitized effective configuration section before numeric results, and benchmark comparisons must cite matching or intentionally varied config fields. |

## Coverage Summary

- Active requirements: 4
- Mapped to slices: 4
- Validated: 0
- Unmapped active requirements: 0
