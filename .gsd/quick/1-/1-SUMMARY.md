# Quick Task: изучи проект

**Date:** 2026-05-19
**Branch:** master

## What Changed
- Изучен проект `fd`: высокопроизводительный сервис эмбеддингов с Go API, двухуровневым кешем L1 `sync.Map` + L2 Redis и TEI-инференсом `deepvk/USER-bge-m3`.
- Зафиксирована архитектура: `api/main.go` собирает зависимости, `/health`, `/v1/embeddings` и `/embeddings/batch`; обработчики используют `cache.TieredCache`; TEI клиент вызывает OpenAI-compatible `/embeddings`; Redis хранит бинарные эмбеддинги с dimension-aware ключами.
- Отмечены ключевые режимы: 1024d для nodes, 512d для edges, OpenAI-compatible ответ для `/v1/embeddings`, base64/float encoding для batch endpoint.

## Files Modified
- `.gsd/quick/1-/1-SUMMARY.md`

## Verification
- Изучены `README.md`, `CHANGELOG.md`, Docker Compose конфигурация и основные Go-файлы в `api/`.
- `cd api && go test ./... -short` — passed (34 tests across 4 packages).
