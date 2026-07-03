# M055-ucu1jl: Phase 1b miss-coalescing

**Vision:** Phase 1b из Issue #9: cross-request miss-coalescing через CoalescingEmbedder. Concurrent /v1/embeddings запросы в time-window (5ms) сливаются в один downstream TEI call. Опционально через FD_COALESCE_ENABLED (default false), backward compatible. Синтетический baseline на 44-ФЗ corpus доказывает 96-97% reduction в TEI calls при burst.

## Slices

- [x] **S01: CoalescingEmbedder + 44-FZ baseline** `risk:low` `depends:[]`
  > After this: TestCoalescingBaseline44FZProof: 50 concurrent goroutines coalesced в 2 TEI calls vs 50 baseline. p95 latency 110ms→11.5ms.

## Boundary Map

Not provided.
