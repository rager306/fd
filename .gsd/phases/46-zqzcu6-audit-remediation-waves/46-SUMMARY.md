---
id: M046-zqzcu6
title: "Audit remediation waves"
status: complete
completed_at: 2026-06-15T06:08:31.535Z
key_decisions:
  - Use wave-based remediation instead of one broad PR.
  - Keep health/readiness endpoints unauthenticated while protected endpoints fail closed.
  - Use handler-specific validation for batch JSON shapes with shared generic guardrails.
  - Use single lock-owned LocalCache state with explicit Close.
  - Close residual P1 findings in M046 and defer only P2/P3 items with explicit rationale.
  - Resolve false-positive GSD browser validation gate through the official verdict override path, not by manually editing DB.
key_files:
  - benchmark-results/m046-s01-audit-validation.md
  - benchmark-results/m046-s02-batch-guardrails.md
  - benchmark-results/m046-s03-batch-backend-chunking.md
  - benchmark-results/m046-s04-exposure-posture.md
  - benchmark-results/m046-s05-localcache-correctness.md
  - benchmark-results/m046-s06-audit-closure.md
  - documents/issue-3-audit-remediation-plan-m046.md
  - documents/issue-3-current-m046.md
  - api/handlers/embeddings.go
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/cache/local.go
  - api/middleware/auth.go
  - api/middleware/ratelimit.go
  - api/main.go
  - .gsd/milestones/M046-zqzcu6/M046-zqzcu6-VALIDATION.md
lessons_learned:
  - Audit findings should be verified against current code before fixing.
  - DoS-class remediation can use static/test proof instead of abusive load tests.
  - GSD UAT content must use the exact `## UAT Type` template form.
  - Residual lower-priority findings need explicit accepted/deferred rationale to avoid silent audit drift.
  - GSD's browser evidence gate can false-positive on backend audit language; use the official verdict override path when no browser-observable acceptance criteria exist.
---

# M046-zqzcu6: Audit remediation waves

**M046 turned issue #3 into verified, risk-ordered remediation: all P0/P1 findings are fixed and P2/P3 residuals are triaged.**

## What Happened

M046 began by validating GitHub issue #3 instead of blindly applying a broad audit patch. S01 confirmed the P0/P1 clusters and root assumptions. S02 hardened `/v1/batch` and `/embeddings/batch` with body, rate-limit, lifecycle, and input guardrails. S03 removed batch endpoint backend work amplification through bounded cache-miss chunking. S04 changed protected endpoint exposure to fail closed without `FD_API_KEY`, protected `/metrics`, disabled trusted forwarded headers by default, and bounded rate limiter state. S05 made LocalCache deterministic by replacing sync.Map plus a separate size counter with one mutex-owned map, derived size, idempotent Close, and shutdown integration. S06 fixed residual P1 #6 with batched cache peeks/Redis MGET for `/v1/embeddings`, fixed residual P1 #9 by registering canonical `method_not_allowed`, and produced a 32-finding closure matrix. Final gates passed: `go test ./...` 284 tests, cache race 9 tests, lint 0 issues, govulncheck 0 reachable vulnerabilities. The validation blocker was resolved via GSD's official verdict override path after confirming browser-gate was a false positive for backend audit artifacts.

## Success Criteria Results

- ✅ P0/P1 findings verified before fixes: S01 evidence.
- ✅ Batch endpoint abuse paths fixed: S02/S03 evidence.
- ✅ Default exposure/auth behavior fail-closed: S04 evidence.
- ✅ LocalCache correctness fixed with race evidence: S05 evidence.
- ✅ Mandatory gates passed: S06 final gates.
- ✅ Residual P2/P3 findings triaged: S06 closure matrix.
- ✅ Milestone validation verdict: pass after official GSD verdict override for false-positive browser gate.

## Definition of Done Results

- ✅ All six planned slices complete.
- ✅ Each slice has SUMMARY and UAT artifacts.
- ✅ Requirements R029-R032 validated.
- ✅ Final full tests, race tests, lint, and govulncheck passed.
- ✅ Closure matrix covers all 32 issue #3 findings.
- ✅ Validation artifact is `verdict: pass`.
- ✅ Commits are local only; no push performed.

## Requirement Outcomes

- R029 validated: batch endpoint guardrails and backend shaping.
- R030 validated: exposure posture fail-closed.
- R031 validated: LocalCache correctness/lifecycle.
- R032 validated: batched cache peeks for `/v1/embeddings`.

## Deviations

S06 additionally fixed P1 #9 after recheck found it still live. The initial milestone validation was auto-downgraded by a false-positive browser evidence gate; it was resolved through GSD's official pass verdict override path with the false-positive rationale recorded.

## Follow-ups

Future milestone candidates: dependency resilience (#11/#12/#14), API contract polish (#22/#24), and maintainability cleanup (#15/#19/#27/#28/#30/#31/#32).
