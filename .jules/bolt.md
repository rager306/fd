## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-10-04 - Optimize Redis cache key allocation
**Learning:** In Go, `strconv.Itoa` only caches integers 0-99. For larger numbers (like 512 and 1024), it results in heap allocations on the hot path (like cache key generation).
**Action:** Hardcode string constants for expected common inputs to enable the compiler to perform a single allocation for the whole string concatenation.
