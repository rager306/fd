# S01: S01

**Goal:** Verify claims about Go BGE-M3/ONNX embedding libraries and adjacent alternatives from current sources.
**Demo:** After this, the suggested Go embedding libraries are classified as viable, risky, or not relevant based on source evidence.

## Must-Haves

- `go-bge-m3-embed` is checked from source docs/repo.
- `all-minilm-l6-v2-go` or equivalent MiniLM Go/ONNX path is checked.
- General ONNX Runtime Go bindings path is checked.
- Findings distinguish facts from assumptions.

## Proof Level

- This slice proves: web/source research

## Integration Closure

Research inputs become trustworthy before any runtime design decision.

## Verification

- Records source links, maintenance signals, runtime requirements, and unknowns.

## Tasks

- [x] **T01: Verified `go-bge-m3-embed` exists and is relevant, but carries artifact/deployment and maturity risk.** `est:small`
  Search and read current source information for `github.com/Dsouza10082/go-bge-m3-embed`, including README, dependencies, model/runtime claims, and maintenance signals.
  - Verify: Source URLs read and findings recorded with facts vs unknowns.

- [x] **T02: Verified MiniLM Go ONNX path as a reference only, not a model replacement under the Russian/legal constraint.** `est:small`
  Search and read current source information for MiniLM Go ONNX examples such as `github.com/clems4ever/all-minilm-l6-v2-go`, and classify relevance to BGE-M3 replacement.
  - Verify: Source URLs read and relevance classification recorded.

- [x] **T03: Research model-preserving ONNX Runtime Go path** `est:medium`
  Research ONNX Runtime Go bindings and BGE-M3 ONNX export/runtime considerations, including native library deployment, quantization, CPU tuning, and whether model-preserving BGE-M3 ONNX artifacts can be produced and validated against the current fd runtime.
  - Verify: Source URLs read and operational implications recorded.

- [x] **T04: Define Russian legal corpus benchmark gate** `est:small`
  Define minimum Russian legal corpus benchmark requirements for any model-changing optimization path: corpus shape, query/relevance judgments, metrics, and baseline comparison against current model.
  - Verify: Benchmark gate recorded and referenced by later recommendation.
