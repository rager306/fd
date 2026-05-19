---
id: M001-h8xt3d
title: "Review remediation"
status: complete
completed_at: 2026-05-19T07:09:21.750Z
key_decisions:
  - Use dimension-aware keys in every cache coordination layer that can store multiple embedding dimensions.
  - Keep Redis write failures best-effort after successful embedding validation in TieredCache.
  - Use handler dependency interfaces so tests exercise production handlers.
  - Keep Redis host port exposure in docker-compose.override.yaml for local dev while base compose stays safer.
key_files:
  - api/cache/tiered.go
  - api/cache/redis.go
  - api/cache/local.go
  - api/cache/redis_binary_test.go
  - api/cache/tiered_cache_test.go
  - api/cache/local_test.go
  - api/handlers/embeddings.go
  - api/handlers/batch.go
  - api/handlers/embeddings_integration_test.go
  - api/main.go
  - api/Dockerfile
  - docker-compose.yaml
  - docker-compose.override.yaml
  - README.md
  - .gsd/milestones/M001-h8xt3d/M001-h8xt3d-VALIDATION.md
lessons_learned:
  - The active GitNexus index can be too broad to analyze a nested repo; local text-search blast radius should be documented when graph tooling is unavailable.
  - Tests that copy handler logic can hide production regressions; small interfaces provide better test seams.
---

# M001-h8xt3d: Review remediation

**Review remediation completed: cache correctness, API validation, LocalCache semantics, and runtime config hardening are fixed and verified.**

## What Happened

Completed the review remediation milestone across four slices. S01 fixed cache dimension isolation and short-vector panic safety. S02 made batch validation strict and moved tests onto production handlers through minimal dependency interfaces. S03 made LocalCache overwrite, TTL refresh, and maxSize behavior real and tested. S04 hardened Docker/runtime configuration by installing curl for healthchecks, aligning REDIS_HOST/PORT config, moving Redis host exposure to the local override, and documenting PORT. Final validation passed with Compose config checks and the full short Go test suite.

## Success Criteria Results

- Cache correctness: passed via S01 targeted tests.
- API validation: passed via S02 handler tests.
- Runtime hardening: passed via S04 Compose config checks.
- Full regression: passed via final `cd api && go test ./... -short`.

## Definition of Done Results

- [x] High-risk cache correctness bugs are fixed and covered by tests.
- [x] API validation is consistent across endpoints.
- [x] Docker health and Redis exposure defaults are safer.
- [x] Handler tests exercise production handlers instead of copied logic where feasible.
- [x] All Go tests pass with `cd api && go test ./... -short`.
- [x] Changes are committed in logical commits with GSD artifacts updated.

## Requirement Outcomes

All findings from the project review were either fixed or documented as non-blocking follow-up.

Evidence: M001-h8xt3d-VALIDATION.md, S01-S04 summaries/UAT files, and final verification output: `base redis port not exposed`, `Go test: 46 passed in 4 packages`.

## Deviations

GitNexus impact/change detection tools could not analyze /root/fd because the active index is scoped to /root and not this repository; text-search blast-radius and git diffs were used instead. Docker image build was not run; Docker Compose config and Go test verification were used.

## Follow-ups

Optional future cleanup: remove obsolete top-level `version` from docker-compose.yaml to silence Compose warning.
