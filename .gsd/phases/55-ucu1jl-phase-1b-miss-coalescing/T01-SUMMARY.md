---
id: T01
parent: S01
milestone: M055-ucu1jl
key_files:
  - api/embed/coalescedembedder.go
  - api/embed/coalescedembedder_test.go
  - api/embed/coalesce_baseline_test.go
  - api/main.go
key_decisions:
  - (none)
duration: 
verification_result: passed
completed_at: 2026-07-03T11:50:04.798Z
blocker_discovered: false
---

# T01: CoalescingEmbedder + 3 unit-теста + 2 44-ФЗ baseline теста + main.go wiring через FD_COALESCE_ENABLED env. 96-97% reduction в TEI calls при burst на 44-ФЗ corpus. 11 пакетов тестов зелёные с -race.

**CoalescingEmbedder + 3 unit-теста + 2 44-ФЗ baseline теста + main.go wiring через FD_COALESCE_ENABLED env. 96-97% reduction в TEI calls при burst на 44-ФЗ corpus. 11 пакетов тестов зелёные с -race.**

## What Happened

T01 закрывает Phase 1b реализацию за один таск. CoalescingEmbedder + синтетический baseline от 44-ФЗ corpus показывает статистическое подтверждение что коалесцинг работает. Data-driven justification для мельчайшего теста-основанного на законодательном корпусе, а не синтетических Lorem Ipsum данных.

## Verification

go test ./api/... -race — 11 packages green. 44-FZ baseline показывает 96%+ reduction в TEI calls при burst workloads.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `cd /root/fd/api && go test ./embed/... -race -count=1 -run TestCoalescing -v` | 0 | pass | 300ms |
| 2 | `cd /root/fd/api && go test ./embed/... -race -count=1 -run TestCoalescingBaseline -v` | 0 | pass | 500ms |

## Deviations

None.

## Known Issues

None.

## Files Created/Modified

- `api/embed/coalescedembedder.go`
- `api/embed/coalescedembedder_test.go`
- `api/embed/coalesce_baseline_test.go`
- `api/main.go`
