---
estimated_steps: 1
estimated_files: 3
skills_used: []
---

# T02: Tier 2 linters enabled and fixed: 17 baseline issues (12 gocritic, 4 gocyclo, 1 unparam) → 0

Fix issues из T01: (a) gocyclo в CreateEmbedding handler — если complexity > 15, refactor: extract helpers (buildEmbeddingRequest, populateResponseFromCache, runChunkedEmbed, applyTruncation). Если после refactor всё ещё > 15, добавить //nolint:gocyclo с explicit comment (handler is inherently multi-stage). (b) durationcheck — fix any int * time.Duration conversions if found. (c) contextcheck — fix any context.Background() in request hot path (сейчас 2 в main.go startup, OK). (d) unparam — fix unused function parameters. (e) nilnil — fix `return nil, nil` patterns if found. (f) gocritic — selective fixes per enabled-tag, not all 30+ checks.

## Inputs

- None specified.

## Expected Output

- `api/handlers/embeddings.go (refactored)`
- `api/embed/codec.go (modified)`

## Verification

golangci-lint run exit 0. CreateEmbedding complexity ≤ 15 (или //nolint:gocyclo с justification). Все M041 acceptance tests pass после refactor.
