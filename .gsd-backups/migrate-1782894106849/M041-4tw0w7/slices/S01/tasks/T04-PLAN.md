---
estimated_steps: 1
estimated_files: 7
skills_used: []
---

# T04: Wire validation в handlers, replace error envelopes, recovery wrapper, 405/404 handlers, encoding_format в /v1/embeddings

(a) api/handlers/embeddings.go: убрать существующие ad-hoc error responses (те что возвращают gin.H{errorKey:...}), использовать errors.WriteError из T02. Подключить validation middleware из T03 в router setup. 405 handler для неправильных methods на /v1/embeddings (T-E-8). Recovery middleware обёрнут чтобы возвращал OpenAI envelope с X-Request-Id (T-E-15). (b) api/handlers/batch.go: то же самое — replace error responses. (c) INVESTIGATE dimensions=512 broken: T01 baseline показал /v1/embeddings {"input":["hello"],"dimensions":512} → 500. Root cause options: (i) TEI model --dtype fp16 не поддерживает 512-dim Matryoshka head, (ii) fd handler truncation bug, (iii) TEI returns 1024 но fd expects 512. Investigate через прямой TEI запрос: curl http://127.0.0.1:30080/embed -d '{"input":"hello","dimensions":512}'. Apply fix: если TEI side, fall back to 1024 с warning (или явная 400 dimensions_mismatch); если fd side, fix handler. (d) MOVE encoding_format codec: extract encodeEmbedding и float32SliceToBytes из api/handlers/batch.go в новый api/embed/codec.go. Добавить EncodingFormat *string в api/embed.EmbeddingsRequest. Validation в T03 принимает encoding_format=float|base64. Handler использует codec для response.

## Inputs

- None specified.

## Expected Output

- `api/handlers/embeddings.go (modified)`
- `api/handlers/batch.go (modified)`
- `api/main.go (modified)`
- `api/embed/types.go (modified, EncodingFormat field)`
- `api/embed/codec.go (new, extracted from batch.go)`
- `api/handlers/embeddings_test.go (modified)`
- `api/handlers/batch_test.go (new)`

## Verification

Integration tests против running fd: T-E-1..T-E-8, T-E-15 все pass. B4 (1MB input) → 413 input_too_long (НЕ timeout 500). B7 ({} missing input) → 400 input_required (НЕ 'unexpected end of JSON'). B9 (100 inputs) → 413 batch_too_large (НЕ 500). B10 (GET /v1/embeddings) → 405 (НЕ 404). Recovery wrapped: forced panic → 500 internal_error с X-Request-Id. dimensions=512: либо fix работает (200), либо явная 400 dimensions_mismatch с понятным message. encoding_format=base64 → 200 с base64 string в response (per T-H-5). encoding_format=float → 200 с float array (default). encoding_format=garbage → 400 dimensions_invalid.
