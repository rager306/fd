---
id: T03
parent: S01
milestone: M053-ckvgjw
key_files:
  - README.md
  - .gsd/runtime/M053-ckvgjw/api-startup-logs.txt
  - .gsd/runtime/M053-ckvgjw/metrics-default-config.txt
key_decisions:
  - FD_CACHE_MAX_SIZE и FD_CACHE_LOCAL_TTL в 'Redis and cache' таблице — одна группа физически
  - Sizing guidance: ~4 KB/entry, 10000 ≈ 40 MB
  - FD_CACHE_LOCAL_TTL = 30s default сохранён — same legacy value, быстрая инвалидация для research workloads
duration: 
verification_result: passed
completed_at: 2026-07-03T10:46:13.058Z
blocker_discovered: false
---

# T03: Docker smoke passed (default config INFO log: cache configured l1_max_entries=10000 l1_ttl=30s). README: FD_CACHE_MAX_SIZE + FD_CACHE_LOCAL_TTL в Configuration, sizing guidance (4 KB/entry, 10000 ≈ 40 MB). Evidence в .gsd/runtime/M053-ckvgjw/.

**Docker smoke passed (default config INFO log: cache configured l1_max_entries=10000 l1_ttl=30s). README: FD_CACHE_MAX_SIZE + FD_CACHE_LOCAL_TTL в Configuration, sizing guidance (4 KB/entry, 10000 ≈ 40 MB). Evidence в .gsd/runtime/M053-ckvgjw/.**

## What Happened

T03 закрывает integration verification + documentation:

1. **Docker smoke** — образ пересобран с --no-cache, запущен, health 200. INFO log показывает "cache configured l1_max_entries=10000 l1_ttl=30s l2_namespace=v2" — env vars read correctly, defaults apply.

2. **README.md** обновлён:
   - "Redis and cache" таблица: добавлены FD_CACHE_MAX_SIZE, FD_CACHE_LOCAL_TTL с purpose descriptions
   - sizing guidance: ~4 KB/entry, 10000 ≈ 40 MB

3. **Evidence**:
   - `.gsd/runtime/M053-ckvgjw/api-startup-logs.txt` — startup INFO logs
   - `.gsd/runtime/M053-ckvgjw/metrics-default-config.txt` — /metrics after startup

Phase 1a ready for production: операторы могут увеличить L1 capacity через FD_CACHE_MAX_SIZE без recompiling. Default behavior идентичен legacy — production deploy без env vars не ломается.

## Verification

go test ./api/... -race — 10 packages green. Docker smoke: /health 200, INFO log "cache configured l1_max_entries=10000 l1_ttl=30s". README: FD_CACHE_MAX_SIZE, FD_CACHE_LOCAL_TTL in Configuration table + sizing formula. Evidence в .gsd/runtime/M053-ckvgjw/.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd && docker compose build --no-cache api && docker compose up -d api` | 0 | pass | 120000ms |
| 2 | `docker logs fd_api --since 60s | grep 'cache configured'` | 0 | pass | 200ms |
| 3 | `curl -s http://localhost:8000/health | jq '.status'` | 0 | pass | 80ms |

## Deviations

Phase 1b deferred into metrics baseline — noted in README's observability section. Phase 1c as TEI tuning guide added to Operations → TEI startup section.

## Known Issues

fd_cache_entries{tier=l1} не экспонируется в /metrics пока нет SetRuntimeObservers на новом образе (известный артефакт М052). Информация о L1 entries видна в startup log.

## Files Created/Modified

- `README.md`
- `.gsd/runtime/M053-ckvgjw/api-startup-logs.txt`
- `.gsd/runtime/M053-ckvgjw/metrics-default-config.txt`
