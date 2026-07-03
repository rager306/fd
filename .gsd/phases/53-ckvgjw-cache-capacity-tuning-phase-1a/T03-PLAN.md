---
estimated_steps: 19
estimated_files: 2
skills_used: []
---

# T03: Docker smoke + README cache tuning documentation

Финальная волна: docker smoke + README + cleanup.

1. **Docker smoke**:
   - Пересобрать образ `docker compose build api` (--no-cache для гарантии свежего бинарника)
   - Запустить с custom env: `FD_CACHE_MAX_SIZE=20 docker compose up -d api`
   - Заполнить >20 unique embeddings через `/v1/embeddings`
   - Проверить `/metrics`: `fd_cache_entries{tier="l1"}` ≤ 20, `fd_cache_evictions_total` растёт
   - Evidence log сохранить в `.gsd/runtime/M053-ckvgjw/`

2. **README.md update**:
   - Configuration section: добавить FD_CACHE_MAX_SIZE и FD_CACHE_LOCAL_TTL в "Redis and cache" table
   - New subsection "Cache tuning" под Operations — sizing guidance, eviction tradeoffs, Phase 1a context per Issue #9
   - Phase 1c guidance: brief note в Operations → TEI startup про tokenization-workers tuning

3. **Phase 1b/1c documentation**:
   - README add: "Phase 1b deferred до data-driven decision от Phase 0 metrics" — объяснить criterion (fd_tei_batch_fill_ratio avg <0.3 → justified)
   - README add: TEI tuning guidance — brief note про ORT threads / tokenization-workers для production

4. **Verify**: 
   - go test ./api/... -race -count=1 — green
   - curl /health 200
   - /v1/embeddings smoke pass
   - README renders cleanly

## Inputs

- `README.md`
- `api/main.go (after T01)`

## Expected Output

- `README.md`
- `.gsd/runtime/M053-ckvgjw/`

## Verification

go test ./api/... -race — green. Docker smoke с FD_CACHE_MAX_SIZE=20 показывает eviction под нагрузкой. README renders, .env.example синхронизирован.

## Observability Impact

README documents env knobs + tuning guidance для operators.
