# M048 Issue #7 Closure Matrix

Captured: 2026-06-15

Issue: https://github.com/rager306/fd/issues/7
Input artifact: `documents/issue-7-current-m048.md`

## Verification Summary

Final gates:

```text
go test ./...: 281 passed in 10 packages
golangci-lint: 0 issues
govulncheck: 0 reachable vulnerabilities
```

## Closure Matrix

| Issue #7 finding | Priority | Status | Evidence |
|---|---:|---|---|
| #19 `LRUCache` is dead production code | P3 | Fixed | S01 deleted LRUCache source/tests and replaced the only integration-test scaffold with a LocalCache-backed adapter. Evidence: `benchmark-results/m048-s01-cache-cleanup.md`; R037 validated. |
| #24 `handleBindError` malformed array-element message | P3 | Fixed | S03 emits `input must be an array of strings, got <Value>` for empty `json.UnmarshalTypeError.Field` and adds a regression test. Evidence: `benchmark-results/m048-s03-api-polish-closure.md`; R039 validated. |
| #26 ONNX-only RuntimeHealth fields | P3 | Fixed | S02 removed inactive ONNX-only fields from RuntimeHealth and removed stale ONNX health test. Evidence: `benchmark-results/m048-s02-runtime-contract-cleanup.md`; R038 validated. |
| #27 duplicate hash-truncate helpers | P3 | Fixed | S01 added canonical cache `shortHash` and removed `shortCacheKeyHash`/`shortNamespaceHash`. Evidence: `benchmark-results/m048-s01-cache-cleanup.md`; R037 validated. |
| #28 `envInt` copies | P3 | Fixed | S01 added `internal/envutil` and updated active main/rate-limit parsing to use shared helpers; LRU env parser was removed with LRU. Evidence: `benchmark-results/m048-s01-cache-cleanup.md`; R037 validated. |
| #29 duplicate embed interfaces | P3 | Fixed | S02 added shared `embed.Embedder` and removed handler/lifecycle duplicate interface declarations. Evidence: `benchmark-results/m048-s02-runtime-contract-cleanup.md`; R038 validated. |
| #30 lifecycle default singleton | P3 | Fixed | S02 removed `defaultState`/`DefaultState()` and main now explicitly calls `lifecycle.NewState()`. Evidence: `benchmark-results/m048-s02-runtime-contract-cleanup.md`; R038 validated. |
| #31 `openapi.m()` silently drops non-string keys | P3 | Fixed | S03 changes `openapi.m()` to panic on non-string keys and adds a regression test. Evidence: `benchmark-results/m048-s03-api-polish-closure.md`; R039 validated. |

## Requirement Outcomes

- R037 validated: cache cleanup consolidation for #19/#27/#28.
- R038 validated: runtime contract simplification for #26/#29/#30.
- R039 validated: API polish for #24/#31.

## Notes

- Browser/runtime UI verification is not applicable; issue #7 is backend cleanup and contract polishing.
- No GitHub issue comments or closure actions were performed by this milestone.
