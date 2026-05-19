---
id: T04
parent: S01
milestone: M008-6hnowu
key_files: []
key_decisions: []
duration: 
verification_result: mixed
completed_at: 2026-05-19T17:04:53.046Z
blocker_discovered: false
---

# T04: Defined the Russian legal benchmark gate: no quality-risking optimization ships without Recall/nDCG/MRR on Russian legal corpus versus current baseline.

**Defined the Russian legal benchmark gate: no quality-risking optimization ships without Recall/nDCG/MRR on Russian legal corpus versus current baseline.**

## What Happened

Defined the Russian legal corpus benchmark gate. Any model-changing or quality-risking optimization (model substitution, INT8, sparse/ColBERT ranking changes, tokenizer/pooling change) must be evaluated on a Russian legal retrieval corpus, not only latency. Minimum dataset shape: corpus of legal passages/chunks with stable IDs, query set in Russian legal language, relevance judgments or gold supporting passages per query, chunking/version metadata, and baseline outputs from the current TEI/Candle BGE-M3-style runtime. Minimum metrics: Recall@k, nDCG@k, MRR@k, optionally MAP/Precision@k, plus latency/resource metrics. RusBEIR is relevant because it is an open Russian IR benchmark with 17 datasets and includes neural models such as BGE-M3, but it may not be domain-perfect. Legal retrieval benchmark sources show domain-specific retrieval quality matters and legal tasks can require reasoning/gold support passages. Acceptance should require no meaningful regression versus current model on Russian legal retrieval metrics while improving speed/resource usage enough to justify risk.

## Verification

Read Russian IR benchmark, legal retrieval benchmark, and IR metric sources; converted into minimum project gate.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `Search query: Russian legal retrieval benchmark dataset embeddings information retrieval metrics nDCG MRR Recall legal Russian corpus` | -1 | unknown (coerced from string) | 0ms |
| 2 | `Read evidence for RusBEIR Dialogue 2025 Russian IR benchmark` | -1 | unknown (coerced from string) | 0ms |
| 3 | `Read legal retrieval benchmark evidence from Stanford legal RAG benchmark and Voyage legal embedding writeup` | -1 | unknown (coerced from string) | 0ms |
| 4 | `Read IR metric references for nDCG@k, MRR@k, Recall@k` | -1 | unknown (coerced from string) | 0ms |

## Deviations

None.

## Known Issues

A public Russian legal benchmark may not exactly match the target corpus. RusBEIR provides Russian IR structure and BGE-M3 evidence, while legal-specific public benchmarks are often English or other jurisdictions. For product-specific legal Russian quality, a project-owned evaluation set may still be required.

## Files Created/Modified

None.
