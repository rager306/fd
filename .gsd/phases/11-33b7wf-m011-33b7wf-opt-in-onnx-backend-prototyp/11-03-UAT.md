# S03: Opt in ONNX dense backend prototype — UAT

**Milestone:** M011-33b7wf
**Written:** 2026-05-19T19:27:25.618Z

# S03 UAT — Opt in ONNX dense backend prototype

## Checks

- [x] Go ONNX dependencies compile.
- [x] ONNX Runtime shared library path is explicit.
- [x] Env-gated live ONNX embedder test passes with local artifact.
- [x] `EMBEDDING_BACKEND=onnx` starts local API when all env vars are provided.
- [x] Shared-cache comparison pitfall discovered and corrected with isolated cache namespace.
- [x] Isolated-cache TEI-vs-Go-ONNX comparison artifact exists.
- [x] Comparison fails cosine threshold, so backend is not equivalent yet.
- [x] Tokenizer mismatch evidence captured.
- [x] ONNX local server was cleaned up.

## UAT Result

Blocked with evidence. ONNX Runtime integration is technically viable, but Go tokenizer parity is not solved. Do not benchmark or recommend production integration until tokenizer output matches Hugging Face tokenizer behavior.

