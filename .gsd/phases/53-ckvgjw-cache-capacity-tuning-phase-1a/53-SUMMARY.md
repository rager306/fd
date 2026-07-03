---
id: M053-ckvgjw
title: "Cache capacity tuning (Phase 1a)"
status: complete
completed_at: 2026-07-03T10:47:06.246Z
key_decisions:
  - env-tuning через envutil.PositiveInt + DurationOrDefault — reusable, tested helpers
  - Defaults = hardcoded legacy values — zero risk production deploy без env vars
  - Duration parsing falls back к default на invalid (negative, unparseable)
  - INFO log on startup for operator visibility
  - Phase 1b deferred to Phase 0 metrics baseline (fd_tei_batch_fill_ratio)
key_files:
  - api/main.go
  - api/internal/envutil/duration.go
  - api/internal/envutil/duration_test.go
  - api/main_test.go
  - .env.example
  - README.md
lessons_learned:
  - envutil package as a centralized parsing hub — future env vars should reuse PositiveInt/DurationOrDefault/BoolOrDefault
  - INFO log on subsystem startup gives operators direct visibility without needing to read /metrics
---

# M053-ckvgjw: Cache capacity tuning (Phase 1a)

**Phase 1a из Issue #9 завершён: env-tunable L1 cache capacity (FD_CACHE_MAX_SIZE) и TTL (FD_CACHE_LOCAL_TTL) с backward-compatible defaults. DurationOrDefault helper. 5 новых unit тестов. Docker smoke INFO log подтверждает config. README extended с sizing guidance.**

## What Happened

M053 завершён. Phase 1a cache capacity env-tuning из Issue #9 реализован. Phase 1b/1c deferred до data-driven decision baseline / README guidance.

## Success Criteria Results

7/7 UAT checks PASS. All success criteria met with objective evidence.

## Definition of Done Results

All DOD items met. 3 tasks, 1 slice, env-tunable cache, 5 unit tests, README, evidence.

## Requirement Outcomes

No new requirements — operational tuning per Issue #9.

## Deviations

None significant. Phase 1c as docs-only deferred — README already had enough TEI guidance.

## Follow-ups

Phase 1b (miss-coalescing): после observation period через Phase 0 metrics. Decision criterion: fd_tei_batch_fill_ratio avg <0.3 → оправдано, иначе пропустить.
