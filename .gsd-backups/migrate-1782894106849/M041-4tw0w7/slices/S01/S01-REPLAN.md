# S01 Replan

**Milestone:** M041-4tw0w7
**Slice:** S01
**Blocker Task:** T01
**Created:** 2026-06-13T18:01:10.750Z

## Blocker Description

Recon + baseline (T01) выполнены в planning phase. S01-RECON.md + benchmark-results/fd-v2-baseline-before-m041-s04.md существуют. T01 помечен completed, immutable. Расширяю T04 чтобы покрыл 2 additional bugs обнаруженные baseline'ом: dimensions=512 broken и encoding_format silent fail в /v1/embeddings. User requested: "Сначала S01 руками — зафиксировать validation как можно скорее, чтобы dimensions=512 broken и encoding_format silent fail были разобраны."

## What Changed

T04 expanded: wire validation (как было) + investigation + fix для dimensions=512 broken (root cause: TEI fp16 не поддерживает 512-dim, или fd truncation bug) + перенос encoding_format codec из batch.go в /v1/embeddings. Estimate 3h → 6h. T02, T03, T05 без изменений.
