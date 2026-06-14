---
estimated_steps: 1
estimated_files: 2
skills_used: []
---

# T05: Regression suite: M041 acceptance в async mode

tools/verify_fd_v2_contract.py (M041 deliverable, deferred) запустить с FD_ASYNC_CHUNKS=true. Все 45 test cases должны pass. Любой failure — bug regression в S02. Проверить особенно: encoding_format=base64 с batch=128 (4 chunks, 128 embeddings in base64), dimensions=512, validation envelope (413, 400 etc) — всё должно работать identical к sync mode.

## Inputs

- None specified.

## Expected Output

- `tools/verify_fd_v2_contract.py (M041 deliverable, run with FD_ASYNC_CHUNKS=true)`
- `tests/integration/fd_v2_async_test.go`

## Verification

go test ./tests/integration/... -run TestFdV2AsyncMode: все M041 acceptance pass в async mode. 0 regressions.
