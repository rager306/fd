## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-06-13 - O(N) Cache Lookups vs Bulk L1/L2 Lookups
**Learning:** Per-item sequential `cache.GetIfPresent` lookups inside batch loops cause high mutex contention on local L1 caches and lead to O(N) single-item Redis round-trips on L2 misses.
**Action:** When evaluating arrays or chunks of items against a Tiered or Local Cache, always utilize or create bulk `GetManyIfPresent` methods. This consolidates lock acquisition overhead and enables L2 layers like Redis to process misses in a single `MGET` pipeline, yielding measurable backend performance gains.
