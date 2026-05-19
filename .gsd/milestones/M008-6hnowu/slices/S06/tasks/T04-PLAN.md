---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T04: Verify BGE-M3 ONNX multi-output implications

Verify BGE-M3 ONNX output claims: dense, sparse, and ColBERT outputs, how they differ from fd's current dense embedding API/cache, and whether fd should ignore, expose, or separately store non-dense outputs.

## Inputs

- `api/embed/types.go`
- `api/cache/redis.go`
- `.gsd/REQUIREMENTS.md`
- `.gsd/DECISIONS.md`

## Expected Output

- `S06 T04 summary`

## Verification

Output-shape implications and compatibility decisions recorded.

## Observability Impact

Maps multi-output model artifacts to fd API/cache compatibility and future search capabilities.
