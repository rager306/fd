# S03: Runtime cache validation — UAT

**Milestone:** M003-xx4yc3
**Written:** 2026-05-19T08:19:13.818Z

# UAT: S03 Runtime cache validation

## Verification performed

- 1024d request after Redis flush created `embed:cache:v2:*:d1024` with STRLEN `4098`.
- Same text at 1024d and 512d created two keys with same hash and suffixes `:d1024` / `:d512`; sizes `4098` / `2050`.
- Cold request: `0.28s`; warm request: `0.00s`; cache-miss log count increased only for cold.
- After API restart, same cached request completed in `0.01s`; Redis key remained; cache-miss log count unchanged.

