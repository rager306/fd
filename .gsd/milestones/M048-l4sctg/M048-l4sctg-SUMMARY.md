---
id: M048-l4sctg
title: "Audit cleanup tail"
status: complete
completed_at: 2026-06-15T11:35:45.174Z
key_decisions:
  - Delete LRUCache rather than keep a dormant reservation because production code uses active LocalCache/TieredCache paths.
  - Use `internal/envutil.Int` and `PositiveInt` to preserve zero-allowed vs positive-only config semantics.
  - Centralize the inference interface in package `embed`, where the concrete TEI client lives.
  - Remove ONNX-only health fields because TEI-only is the active product path.
  - Make `openapi.m()` panic on developer misuse rather than silently dropping schema fields.
key_files:
  - api/cache/hash.go
  - api/internal/envutil/int.go
  - api/embed/types.go
  - api/handlers/health.go
  - api/lifecycle/state.go
  - api/middleware/validation.go
  - api/openapi/spec.go
  - benchmark-results/m048-issue-7-closure.md
lessons_learned:
  - Package-local helpers may not be indexed by GitNexus; use targeted tests/static proof when impact lookup returns UNKNOWN.
  - Cleanup slices benefit from closure matrices tying every reported finding to a file-level proof and requirement outcome.
---

# M048-l4sctg: Audit cleanup tail

**M048 closes GitHub issue #7’s eight P3 cleanup findings with dead-code removal, helper consolidation, runtime contract simplification, and API polish.**

## What Happened

M048 resolved GitHub issue #7 across three small cleanup slices. S01 removed dead LRUCache production code, replaced the only LRU-based integration test scaffold with an active LocalCache-backed adapter, unified duplicate cache hash helpers, and centralized active env integer parsing. S02 simplified runtime contracts by removing inactive ONNX-only health fields, centralizing the embedding interface in package embed, and removing the lifecycle default singleton. S03 polished API helper behavior by fixing the malformed validation message for non-string array input and making `openapi.m()` fail loudly on non-string keys. The milestone also wrote a full issue #7 closure matrix and validated requirements R037-R039. No GitHub outward action was performed.

## Success Criteria Results

- ✅ Issue #7 findings #19/#24/#26/#27/#28/#29/#30/#31 revalidated and closed in `benchmark-results/m048-issue-7-closure.md`.
- ✅ Dead LRU cache production code removed; tests use active LocalCache/TieredCache paths.
- ✅ Duplicate hash helpers and active env integer parsing copies unified.
- ✅ Runtime health and embedding/warmup contracts reflect active TEI-only path.
- ✅ Validation and OpenAPI helper behavior now fails clearly for malformed inputs/developer misuse.
- ✅ Full tests, lint, govulncheck, artifact UAT, and milestone validation passed.

## Definition of Done Results

- ✅ All planned slices complete: S01, S02, S03.
- ✅ Requirements validated: R037, R038, R039.
- ✅ Final `go test ./...`: 281 passed in 10 packages.
- ✅ Final golangci-lint: 0 issues.
- ✅ Final govulncheck: 0 reachable vulnerabilities.
- ✅ Milestone validation verdict: pass.

## Requirement Outcomes

| Requirement | Status | Proof |
|---|---|---|
| R037 | Validated | S01 removed dead/duplicate cache cleanup tail and passed tests/static proof. |
| R038 | Validated | S02 simplified runtime/lifecycle contracts and passed tests/static proof. |
| R039 | Validated | S03 fixed validation/OpenAPI helper behavior and passed tests/static proof. |

## Deviations

S03 final lint required adding a package comment to the new `internal/envutil` package. No scope deviations.

## Follow-ups

Optional outward action only after explicit user confirmation: push local M048 commits and close/comment on GitHub issue #7 with the closure matrix.
