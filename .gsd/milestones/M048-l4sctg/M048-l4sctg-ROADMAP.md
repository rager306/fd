# M048-l4sctg: Audit cleanup tail

**Vision:** Resolve GitHub issue #7 by cleaning the remaining P3 maintainability tail without changing fd's active TEI runtime behavior. Each cleanup removes misleading or duplicated surfaces only after tests/static checks prove the current behavior and the post-cleanup contract.

## Success Criteria

- Issue #7 findings #19, #24, #26, #27, #28, #29, #30, and #31 are revalidated against current code before changes.
- Dead LRU cache production code is removed or explicitly justified, and tests use active LocalCache/TieredCache paths.
- Duplicate cache hash helpers and active env integer parsing copies are unified without changing runtime defaults.
- Runtime health and embedding/warmup contracts reflect the active TEI-only product path.
- Validation and OpenAPI helper behavior fails clearly for malformed inputs.
- Full Go tests, lint, govulncheck, artifact UAT, and milestone validation pass.

## Slices

- [x] **S01: Cache cleanup consolidation** `risk:medium` `depends:[]`
  > After this: Dead LRU cache code and duplicate cache helper/env parsing surfaces are removed or unified with proof that active cache behavior remains green.

- [x] **S02: Runtime contract simplification** `risk:medium` `depends:[S01]`
  > After this: Health and lifecycle contracts expose only active TEI runtime surfaces and one embedding interface contract.

- [x] **S03: API polish and closure** `risk:medium` `depends:[S01,S02]`
  > After this: Validation messages and OpenAPI helper failures are clear, then issue #7 closes with final gates and closure matrix.

## Boundary Map

Not provided.
