---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T02: FD_BACKEND env selection в main.go

api/main.go: parse FD_BACKEND env ("tei" | "onnx", default "tei"). Если "onnx" — instantize ONNXEmbedder вместо TEIClient, передать в handlers.NewEmbeddingsHandler. Должно работать БЕЗ перекомпиляции если бинарь уже собран с -tags onnx. Если FD_BACKEND=onnx но бинарь без onnx tag — fail fast с clear error message. Также env FD_ONNX_ARTIFACT_MANIFEST, FD_ONNX_RUNTIME_LIBRARY, FD_ONNX_TOKENIZER_PATH уже читаются (M008). Создать /info endpoint data update: добавить backend: "tei"|"onnx" поле.

## Inputs

- None specified.

## Expected Output

- `api/main.go (modified)`
- `api/handlers/observability.go (modified)`

## Verification

FD_BACKEND=onnx → ONNXEmbedder used. FD_BACKEND=tei (default) → TEIClient used. /info endpoint показывает backend field. Unit test: выбор backend based on env.
