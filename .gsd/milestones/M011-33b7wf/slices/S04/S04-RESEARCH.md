# S04 — Research

**Date:** 2026-05-20

## Summary

M011 should close as an evidence-backed blocked prototype, not as a production-ready ONNX backend. The milestone successfully established the artifact/config contract, kept TEI as the default runtime, added explicit ONNX startup validation, and proved that a Go ONNX Runtime CPU EP path can load and execute the local `deepvk/USER-bge-m3` dense FP32 artifact. The remaining blocker is semantic: the Go tokenizer path used in S03 does not match Hugging Face tokenization for the same Russian legal probe, and the resulting embeddings fail the existing cosine equivalence threshold.

The isolated-cache comparison artifact `benchmark-results/fd-go-onnx-m011-s03.txt` is the decisive evidence. With `EMBEDDING_CACHE_VERSION=m011-onnx`, the Go ONNX backend returned normalized 1024-dimensional vectors, but TEI-vs-Go-ONNX cosine values ranged from `0.98266755` to `0.99713198`, below the `0.999` gate. A token-ID probe then showed the likely root cause: Hugging Face Python tokenization produced 21 IDs for a Russian labor-law probe, while `sugarme/tokenizer` produced 27 IDs with divergent IDs.

Do not benchmark ONNX throughput yet. Speed numbers would be misleading until preprocessing parity is solved, because the system would be measuring a different embedding function. The next milestone should focus on tokenizer parity first, then rerun equivalence and only then run TEI-vs-ONNX performance comparisons.

## Recommendation

Close M011 with the recommendation: **continue ONNX work only through a tokenizer parity milestone; keep TEI as default and do not make production-readiness claims.**

The next milestone should build a small tokenizer parity harness that compares Hugging Face Python tokenizer output against candidate Go tokenization approaches on the fixed Russian/legal probe set. Acceptance should be token-level equality for `input_ids` and `attention_mask` before ONNX embedding cosine is considered meaningful. If a pure-Go tokenizer cannot match within a short spike, the project should either defer Go ONNX or evaluate a narrow tokenizer sidecar/FFI path with explicit operational tradeoffs.

## Implementation Landscape

### Key Files

- `docs/onnx-artifacts/user-bge-m3-dense-fp32.json` — tracked manifest/checksum contract for the local ignored ONNX artifact.
- `api/embed/onnx_manifest.go` — validates manifest metadata, checksum, dimensions, output name, and resolves repo-root-relative artifact paths.
- `api/embed/onnx.go` — real Go ONNX embedder prototype using explicit shared library, tokenizer path, and dynamic ONNX Runtime session.
- `api/main.go` — backend selection seam; TEI remains default, ONNX requires explicit env vars.
- `benchmark-results/fd-go-onnx-m011-s03.txt` — failed isolated-cache TEI-vs-Go-ONNX comparison; current blocker evidence.
- `.gsd/milestones/M011-33b7wf/slices/S03/S03-SUMMARY.md` — detailed S03 outcome and lessons.

### Build Order

1. Close M011 without further ONNX performance work.
2. Plan tokenizer parity as the next milestone.
3. In that milestone, first produce a Python HF tokenizer baseline artifact with labels, char counts, `input_ids`/`attention_mask` hashes, lengths, and no raw probe text.
4. Compare Go tokenizer candidates against the same baseline.
5. Only after token parity passes, rerun Go ONNX vs TEI cosine using an isolated cache namespace.
6. Only after cosine equivalence passes, benchmark latency/throughput/startup/memory.

### Verification Approach

For M011 closure:

- `cd api && go test ./... -short`
- `cd api && go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run --config ../.golangci.yml ./...`
- `docker compose config`
- `curl -fsS http://localhost:8000/health`
- parse `docs/onnx-artifacts/user-bge-m3-dense-fp32.json`
- verify `benchmark-results/fd-go-onnx-m011-s03.txt` exists, records `raw_probe_texts_logged=false`, and has verdict `FAIL`
- verify no large `.onnx` artifacts are tracked
- run GitNexus change detection before commit

## Constraints

- TEI remains the production/default runtime.
- ONNX remains opt-in and must never silently fall back to TEI.
- The large ONNX artifact stays ignored/local; git tracks only manifest and evidence.
- Raw benchmark/comparator probe texts must not be printed in artifacts.
- Redis cache namespaces must be isolated or flushed when comparing backend implementations; otherwise TEI cache hits can mask ONNX behavior.
- The current live ONNX runtime test uses a uv-cache `libonnxruntime.so.1.26.0`, which is not a production artifact distribution strategy.

## Common Pitfalls

- **Benchmarking before equivalence** — latency numbers are invalid if tokenizer parity fails; solve token IDs first.
- **Shared Redis cache masking** — using the same cache namespace can return TEI vectors from Redis even when the ONNX server is running; isolate `EMBEDDING_CACHE_VERSION` for backend comparisons.
- **Treating ONNX load success as correctness** — ONNX Runtime can load and execute while still producing semantically different vectors due to preprocessing differences.
- **Leaking raw legal probe text** — artifacts should use labels, lengths, hashes, and metrics only.

## Open Risks

- A pure-Go tokenizer library may not exactly reproduce Hugging Face tokenizer behavior for this model/tokenizer JSON.
- A tokenizer sidecar could solve parity but adds operational complexity and may erase some expected performance benefit.
- Stable ONNX Runtime shared library packaging remains unresolved.
- Larger Russian/legal corpus validation is still required even after fixed-probe cosine passes.

## Sources

- S03 failed isolated-cache comparison: `benchmark-results/fd-go-onnx-m011-s03.txt`
- S03 completion summary: `.gsd/milestones/M011-33b7wf/slices/S03/S03-SUMMARY.md`
- M010 ONNX feasibility evidence: `benchmark-results/fd-onnx-fp32-m010-s03.txt`
