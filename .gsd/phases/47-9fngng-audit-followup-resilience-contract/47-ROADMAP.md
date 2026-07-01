# M047-9fngng: Audit followup resilience contract

**Vision:** Resolve GitHub issue #6 by turning the remaining G4 and G5 audit follow-ups into tested reliability and contract fixes. The milestone preserves fd's same-host TEI-only runtime assumptions while making dependency failure, warmup failure, shutdown failure, env parsing, and error-code contracts explicit and verifiable.

## Success Criteria

- Issue #6 findings #11, #14, #13, #32, #25, and #15 are revalidated against current code before fixes land.
- `getEnvInt` cannot overflow or accept negative values in a way that disables capacity protection.
- Fatal listener errors no longer call `os.Exit(1)` from the listener goroutine and use `errors.Is` for `http.ErrServerClosed`.
- TEI transient failures have bounded retry/backoff and repeated outage fast-fail behavior with tests.
- Warmup transient failures retry within bounded policy and readiness clears after a later success.
- Dead or reserved error codes are resolved by removal or explicit emitter evidence.
- Full Go tests, focused package tests, lint, govulncheck, and artifact UAT pass.

## Slices

- [x] **S01: Contract cleanup baseline** `risk:medium` `depends:[]`
  > After this: Issue #6 small contract findings have red tests and fixes for env integer parsing and error code registry policy.

- [x] **S02: Graceful server error path** `risk:high` `depends:[S01]`
  > After this: Fatal listener errors enter controlled shutdown handling instead of exiting inside the listener goroutine.

- [x] **S03: TEI retry and fast fail** `risk:high` `depends:[S01]`
  > After this: TEI transient dependency failures retry within a bounded policy and repeated outages fail quickly with a clear error.

- [x] **S04: Warmup retry and closure** `risk:medium` `depends:[S02,S03]`
  > After this: Warmup transient failures retry and the milestone closes issue #6 with final gates and closure matrix.

## Boundary Map

Not provided.
