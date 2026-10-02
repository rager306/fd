## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Redis cache key optimization
**Learning:** Using strconv.Itoa in a highly-frequent hot path (like generating cache keys per embedding lookup) forces heap allocations in Go for integers larger than 99. String literals combined with string concatenation allow the Go compiler to optimize the entire string builder into a single allocation.
**Action:** Use hardcoded fast-paths (e.g. `if dim == 1024`) in cache key builders instead of strconv to avoid allocations and reduce garbage collector overhead without hurting readability.
