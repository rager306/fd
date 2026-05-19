---
id: M005-pyn4x7
title: "Production hardening documentation"
status: complete
completed_at: 2026-05-19T10:45:57.894Z
key_decisions:
  - D001: Keep current measured TEI CPU/Candle fallback runtime; treat ONNX export as future measured optimization requiring A/B benchmark evidence.
  - Benchmark documentation must state destructive/intrusive side effects: Redis FLUSHALL and API restart.
key_files:
  - README.md
  - .gsd/DECISIONS.md
  - .gsd/milestones/M005-pyn4x7/M005-pyn4x7-VALIDATION.md
lessons_learned:
  - Runtime warnings should be classified in docs as correctness bugs, host deployment notes, or measured optimization candidates.
  - Docs for benchmarks must mention side effects and target environment limits, not only command syntax.
---

# M005-pyn4x7: Production hardening documentation

**Documented benchmark workflow, Redis/TEI hardening notes, and recorded the TEI ONNX decision.**

## What Happened

M005 made the runtime validation and optimization knowledge durable in project-facing documentation. README now frames performance numbers as local benchmark snapshots, documents the `uv run --python 3.13 --with requests --with redis python benchmark.py` workflow, and warns that benchmark.py flushes local Redis and restarts the API for Redis L2 diagnostics. README also documents Redis localhost-only exposure, the host-level Redis overcommit warning, LOG_LEVEL/cache debug behavior, and TEI ONNX/Candle fallback implications. GSD decision D001 records that ONNX export is a future measured optimization rather than an immediate correctness requirement. Final verification passed across Compose config, Go tests, benchmark syntax, README snippet checks, and GitNexus low-risk detection.

## Success Criteria Results

- Runtime/benchmark documentation matches actual validated workflow: pass.
- Redis and TEI hardening notes are explicit and actionable: pass.
- No unverified production-hardening code changes introduced: pass.
- GSD artifacts prepared for local commit: pass.

## Definition of Done Results

- README documents uv Python 3.13 benchmark workflow and API restart side effect: met.
- Runtime docs explain Redis localhost binding and overcommit note: met.
- TEI ONNX/Candle runtime decision recorded: met via D001.
- Docker/README guidance matches verified commands/config: met.
- Final verification passes and local commit prepared: verification met; commit follows immediately after DB checkpoint.

## Requirement Outcomes

No formal requirement IDs changed. Launchability and operability are improved because future maintainers can reproduce local benchmarks safely and understand Redis/TEI runtime notes without reading GSD history.

## Deviations

No host sysctl or TEI ONNX artifact changes were applied. M005 intentionally documents those as deployment/future optimization items rather than making unmeasured runtime changes.

## Follow-ups

Push remains pending explicit user confirmation. Future work could add a production deployment checklist or a separate ONNX A/B benchmark milestone if desired.
