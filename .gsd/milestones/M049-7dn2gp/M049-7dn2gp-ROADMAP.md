# M049-7dn2gp: Agent native observability and cache controls

**Vision:** Resolve the actionable solo-operator subset of GitHub issue #8 by adding authenticated cache invalidation plus richer health/metrics context, while deliberately deferring trace multi-tenant hardening and broad config/options extraction until fd needs them.

## Success Criteria

- Issue #8 body is preserved as `documents/issue-8-current-m049.md`.
- AN-A is implemented with tested cache delete/flush primitives and authenticated HTTP invalidation route(s).
- AN-B and AN-C are implemented with health last_error/dependencies/capacity fields plus Prometheus in-flight/cache occupancy signals.
- AN-D is explicitly deferred for solo deployment; AN-E/F are scoped to no broad abstraction unless required by implementation.
- Full Go tests, lint, govulncheck, artifact UAT, live rebuilt-container smoke, and milestone validation pass.

## Slices

- [x] **S01: Cache invalidation controls** `risk:medium` `depends:[]`
  > After this: After this slice, an authenticated operator can flush fd's embedding cache and observe MISS/HIT behavior changing without restarting services.

- [ ] **S02: Health and metrics context** `risk:medium` `depends:[S01]`
  > After this: After this slice, /health and /metrics expose last error, dependency, capacity, and cache occupancy signals agents can inspect.

- [ ] **S03: Solo scope closure and live verification** `risk:low` `depends:[S01,S02]`
  > After this: After this slice, issue #8 has a closure matrix for implemented/deferred items and the rebuilt container passes live smoke.

## Boundary Map

| In scope | Out of scope |
|---|---|
| Authenticated cache delete/flush for fd embedding cache | Multi-tenant trace scoping/admin token split for AN-D |
| Health last_error and dependency reachability/latency | Turning /health into a live embedding inference probe |
| Health capacity and metrics occupancy gauges | Broad pluggable retry/cache policy framework for AN-E/F |
| Closure docs for issue #8 solo scope | OpenAPI /embeddings/batch drift item explicitly excluded by issue #8 |
