---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T01: Audit existing ONNX implementation + build matrix

Прочитать api/embed/onnx.go (уже есть с M008/M019), api/embed/onnx_manifest.go, api/embed/onnx_tokenizer_*.go. Проверить что impl поддерживает batched input (Embed(ctx, texts []string) ([][]float32, error) — это contract). Build matrix: (1) go build . (default, no onnx) — работает, (2) go build -tags onnx . — компилируется, runtime требует libonnxruntime.so. Создать Makefile target `make build-onnx`.

## Inputs

- None specified.

## Expected Output

- `Makefile`
- `api/embed/onnx-audit-m042-s03.md (audit findings)`

## Verification

go build -tags onnx -o /tmp/fd-api-onnx ./api exit 0. Audit doc: ONNXEmbedder implements Embedder interface (batch input supported), handles tokenizer, has artifact manifest validation.
