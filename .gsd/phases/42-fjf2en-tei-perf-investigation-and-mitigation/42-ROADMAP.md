# M042-fjf2en: TEI perf investigation and mitigation

**Vision:** Исследовать почему TEI cold path 6s per batch=32 (queue_time 2.7s) и предоставить операторам choice: (a) async pipeline в fd (env FD_ASYNC_CHUNKS=true) — 4 parallel chunks drop cold 128-batch 25s → 10s; (b) ONNX Go runtime mode (env FD_BACKEND=onnx) — per M019 cold 8ms / warm 1ms, 100-700x faster than TEI, opt-in. Production default остаётся TEI. Все M041 acceptance criteria продолжают pass в обоих режимах.

## Success Criteria

- RCA документ M042 S01 существует и обосновывает mitigation strategy
- Async pipeline (S02) снижает cold path batch=128 с 25s до ≤10s без regression cache hit path
- ONNX mode (S03) opt-in снижает cold path batch=32 с 6s до ≤500ms; default off (TEI per R001)
- All M041 acceptance tests (45 test cases + 10 scenarios) pass в обоих режимах (TEI+sync, TEI+async, ONNX)
- Performance spec targets (docs/fd-v2.md Section 5.4) closing status: T-P-1 ✓ both modes, T-P-2 with ONNX, T-P-3 with ONNX, T-P-5 ✓ both modes (with cache); T-P-2/T-P-3 with TEI still fail (documented как known limitation, requires TEI source change or legal-quality remediation)
- Final artifact benchmark-results/fd-v2-perf-final-m042.md с consolidated table (TEI sync, TEI async, ONNX × batch=1/10/32/128 cold/warm)

## Slices

- [x] **S01: TEI queue_time root cause analysis** `risk:low` `depends:[]`
  > After this: After this, документ docs/te-perf-root-cause-m042.md объясняет почему TEI queue_time=2.7s несмотря на max_concurrent_requests=512, и обосновывает выбор mitigation strategy (async vs ONNX vs both) с evidence.

- [x] **S02: TEI active-path cleanup and safe mitigation** `risk:medium` `depends:[S01]`
  > After this: After this, active build/config/docs no longer present ONNX as a current runtime path; TEI remains the only current backend, and any TEI request-shaping mitigation has fresh safe evidence or is explicitly deferred.

- [x] **S03: ONNX conditional fallback and speed measurement** `risk:high` `depends:[S01,S02]`
  > After this: After this, FD_BACKEND=onnx (requires rebuilding fd binary с -tags onnx) переключает fd на Go ONNX runtime. Per M019: cold path batch=32 ≤500ms, warm path ≤10ms. Default off (TEI остаётся production). Legal quality gate deferred per M015/M016.

## Boundary Map

### S01 → S02

Produces:
- RCA документ обосновывающий async concurrency=4 (TEI max_batch_requests=4)
- Hypotheses resolved для ошибок (что если TEI internal, что если fd client)

Consumes:
- M041 perf baseline (benchmark-results/fd-v2-baseline-before-m041-s04.md)
- M041 handler implementation (api/handlers/embeddings.go chunked loop)

### S02 → S03

Produces:
- Async pipeline baseline (cold path reduced from 25s to ≤10s)
- Bounded concurrency=4 (matches TEI capacity)
- Cache integration unchanged

Consumes:
- S01 RCA verdict
- S01 async concurrency recommendation

### S01 → S03

Produces:
- ONNX as opt-in speed-first alternative
- Cold path ≤500ms (per M019)
- Legal quality gate deferred (per M015/M016)

Consumes:
- M019 ONNX measurements
- M015/M016 legal quality results
- S01 verdict на вопрос "ONNX is viable"
