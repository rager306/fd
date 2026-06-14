---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T02: Live profile TEI: варьировать concurrency, batch size, timing

Active measurement: (1) одновременно 4 curl в parallel с batch=32, измерить max queue_time vs sequential. (2) Одновременно 16 curl с batch=1, измерить queue_time degradation. (3) Sleep 30s, один curl batch=32 — sanity. (4) Перезапустить fd_tei (down/up), один curl batch=32 — измерить true cold start. Цель: понять TEI behavior при разной нагрузке.

## Inputs

- None specified.

## Expected Output

- `tools/profile_tei_concurrency.sh`
- `benchmark-results/te-concurrency-profile-m042-s01.md`

## Verification

Профиль собран, ≥4 сценария, queue_time pattern задокументирован.
