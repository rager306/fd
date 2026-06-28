## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2026-06-28 - Optimize RedisCache key generation
**Learning:** The cache key generation in the Redis cache was doing multiple allocations per call due to `[]byte` casting, `hex.EncodeToString`, and `strconv.Itoa`.
**Action:** Used `unsafe.Slice(unsafe.StringData(text), len(text))` for string-to-byte casting, stack-allocated buffers for `hex.Encode`, and hardcoded `1024`/`512` string fast-paths to reduce allocations in hot path cache key generation.
