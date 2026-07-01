# S05: Research Go vs C vs Rust performance tradeoffs — UAT

**Milestone:** M008-6hnowu
**Written:** 2026-05-19T17:11:43.348Z

# UAT: S05 Research Go vs C vs Rust performance tradeoffs

## Evidence

- T01 mapped fd's bottleneck layers.
- T02 researched Rust maturity and likely gain.
- T03 researched C maturity and likely gain.
- T04 saved rewrite strategy recommendation.

## Acceptance

- No speculative rewrite is recommended.
- Go remains the primary API/cache service.
- Rust is a future sidecar option only after native inference evidence.
- C is limited to narrow FFI/reference usage.

