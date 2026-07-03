---
id: M003-xx4yc3
title: "Runtime validation and performance baseline"
status: complete
completed_at: 2026-05-19T08:28:38.131Z
key_decisions:
  - Use `uv run --python 3.13 --with requests --with redis` for accepted Python benchmark execution.
  - Keep Redis published to host only on 127.0.0.1 so host benchmarks work without exposing Redis externally.
  - Do not mix speculative performance tuning into M003; record baseline first and optimize in a follow-up milestone.
key_files:
  - docker-compose.override.yaml
  - benchmark-results/fd-benchmark-baseline-py313.txt
  - benchmark-results/fd-runtime-stats-logs.txt
  - .gsd/milestones/M003-xx4yc3/M003-xx4yc3-ASSESSMENT.md
  - .gsd/milestones/M003-xx4yc3/M003-xx4yc3-VALIDATION.md
lessons_learned:
  - Compose override files can accidentally expose local infrastructure; bind host-only services explicitly to 127.0.0.1.
  - Benchmark summary code needs validation against raw tables before being treated as authoritative.
  - TEI backend constraints and model artifact format should be measured before tuning API batching.
---

# M003-xx4yc3: Runtime validation and performance baseline

**Validated the fd Docker runtime, fixed Redis host exposure, captured a Python 3.13 benchmark baseline, and documented the next optimization path.**

## What Happened

M003 validated the fd stack in a real Docker runtime. The stack startup path uncovered two issues: a stale local fd_api container caused an initial name conflict, and Redis was exposed on all host interfaces by the override file. The stale container was removed locally and the Redis override was fixed to bind to 127.0.0.1. Live health, embedding, batch, negative validation, and cache tests passed against the running services. Runtime cache validation proved dimension-aware Redis keys and L2 cache survival across API restart. Benchmarking was executed with uv and Python 3.13.12, producing a durable baseline and runtime stats/log artifact. Final gates passed: Docker Compose config, Go tests, and GitNexus low-risk change detection. Remaining work is documented as optimization follow-up rather than runtime correctness blocker.

## Success Criteria Results

- Stack startup/log inspection: pass.
- API/TEI/Redis health: pass.
- Embedding and batch endpoints: pass.
- Negative validation cases: pass.
- Runtime cache behavior: pass.
- Python 3.13 benchmark baseline: pass.
- Resource/log correlation: pass.
- Optimization assessment: pass.
- Final verification: pass.

## Definition of Done Results

- Docker stack exercised: met.
- Logs inspected: met.
- Discovered runtime/security issues fixed or classified: met.
- Live API smoke tests run: met.
- Runtime cache behavior verified: met.
- Benchmark baseline captured under requested uv Python 3.13 runtime: met.
- Findings and optimization plan documented: met.
- Final verification gates passed: met.

## Requirement Outcomes

No formal REQUIREMENTS.md IDs were transitioned, but runtime validation supports the project’s operational capability: API, TEI, Redis, cache, and benchmark paths have live evidence.

## Deviations

Python benchmark was first run with uv default Python before user clarified Python 3.13. It was rerun with `uv run --python 3.13`, and only the Python 3.13 artifact is the accepted baseline. Full TEI cold startup from zero was not measured because TEI was already warmed/running in the environment.

## Follow-ups

Recommended next milestone: fix benchmark.py max-throughput summary selection, add cache metrics/log sampling, extend benchmark modes to isolate L1/L2/cold TEI paths, then evaluate TEI backend/model artifact optimization and batch tuning from the new evidence.
