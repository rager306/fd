## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Redis Cache Key Generation Optimization
**Learning:** `strconv.Itoa` causes an allocation when generating strings dynamically for cache keys, especially in hot paths like Redis cache keys. Fast path static strings (e.g., `if dim == 1024 ...`) skip the `strconv.Itoa` allocation entirely and allow the Go compiler to optimize the entire string concatenation.
**Action:** Replaced `strconv.Itoa` with conditional static strings for common dimensions (1024 and 512) in Redis key generation.
