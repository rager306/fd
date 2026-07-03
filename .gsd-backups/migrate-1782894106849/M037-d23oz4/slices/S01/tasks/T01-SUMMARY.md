---
id: T01
parent: S01
milestone: M037-d23oz4
key_files:
  - api/embed/onnx_test.go
  - api/embed/hf_tokenizer_native_test.go
  - api/main_test.go
  - api/handlers/embeddings_integration_test.go
  - tools/evaluate_legal_retrieval.py
  - benchmark.py
key_decisions: []
duration: 
verification_result: passed
completed_at: 2026-05-21T10:13:08.916Z
blocker_discovered: false
---

# T01: Mapped Go/API/package validation surfaces and the Python-helper boundary.

**Mapped Go/API/package validation surfaces and the Python-helper boundary.**

## What Happened

Inspected current runtime validation surfaces. Go has ONNX constructor/config/health checks, a live local artifact smoke test guarded by env vars, native HF tokenizer parity tests, and API handler integration tests. The legal evaluator and benchmark harness are Python drivers, but they call actual fd API endpoints, so they can be valid target-runtime gates only when pointed at the Go/packaged Go ONNX service with isolated Redis namespaces. This confirms the policy boundary: Python export/provisioning/verifier checks are setup evidence; production acceptance requires target runtime API/package evidence.

## Verification

GitNexus impact for ONNX test was LOW; relevant Go/Python surfaces were read and summarized.

## Verification Evidence

| # | Command | Exit Code | Verdict | Duration |
|---|---------|-----------|---------|----------|
| 1 | `gitnexus_impact(target=TestNewONNXEmbedderRequiresManifestPath, direction=upstream)` | 0 | ✅ pass — LOW risk, no impacted processes | 0ms |
| 2 | `read Go ONNX/config/API tests, legal evaluator, benchmark harness` | 0 | ✅ pass — target runtime surfaces identified | 0ms |

## Deviations

None.

## Known Issues

Existing Go tests include smoke/config/tokenizer checks, but target-runtime acceptance for any new artifact must still include actual Go API legal/performance/package gates, not only Python verifier output.

## Files Created/Modified

- `api/embed/onnx_test.go`
- `api/embed/hf_tokenizer_native_test.go`
- `api/main_test.go`
- `api/handlers/embeddings_integration_test.go`
- `tools/evaluate_legal_retrieval.py`
- `benchmark.py`
