---
id: T01
parent: S04
milestone: M016-pdcjat
key_files:
  - benchmark-results/fd-model-alternatives-m016-s04.txt
key_decisions:
  - Alternative models are research candidates only; none should replace `deepvk/USER-bge-m3` without legal-corpus benchmark evidence.
  - The top alternative-model candidates for future measurement are Qwen3-Embedding-0.6B, multilingual-e5-large/instruct, ai-forever/ru-en-RoSBERTa, d0rj/e5-large-en-ru, and deepvk/USER-base.
  - Rerankers should be considered as a second-stage quality layer, not embedding replacements.
duration: 
verification_result: passed
completed_at: 2026-05-20T05:19:01.353Z
blocker_discovered: false
---

# T01: Researched and ranked alternative embedding/reranker candidates for future Russian legal benchmarks.

**Researched and ranked alternative embedding/reranker candidates for future Russian legal benchmarks.**

## What Happened

Researched current Russian/multilingual embedding alternatives and deployment constraints. Sources covered BGE-M3, Qwen3 Embedding, multilingual E5, Russian E5 derivatives, ai-forever ru-en-RoSBERTa/ruMTEB, deepvk USER-base, and reranker options. The research separates long-context alternatives, Russian-focused candidates, and reference-only small models.

## Verification

Source-backed candidate list with more than 5 candidates was identified and written into the research artifact.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `search_and_read BGE-M3 / Qwen3 / multilingual E5 / Russian embedding candidates` | 0 | ✅ pass — source evidence collected | 0ms |
| 2 | `write benchmark-results/fd-model-alternatives-m016-s04.txt` | 0 | ✅ pass — research artifact drafted | 0ms |

## Deviations

None.

## Known Issues

The research did not download or benchmark alternative models; it ranks candidates from public docs and defines the fair benchmark protocol.

## Files Created/Modified

- `benchmark-results/fd-model-alternatives-m016-s04.txt`
