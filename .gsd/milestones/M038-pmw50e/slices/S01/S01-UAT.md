# S01: Go ONNX runtime smoke proof — UAT

**Milestone:** M038-pmw50e
**Written:** 2026-05-21T10:37:25.162Z

# UAT — M038 S01

Go target-runtime smoke proof passes:

- local artifacts and dependencies verified;
- live Go ONNX embedder test passed;
- Go API started with `EMBEDDING_BACKEND=onnx`, tags `onnx hf_tokenizers`, and namespace `m038-go-onnx-smoke`;
- `/health` reported ONNX backend and verification metadata;
- `/v1/embeddings` returned a 1024-dimensional normalized embedding;
- no raw probe text was recorded;
- server stopped and port 18000 is clean.

