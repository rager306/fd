## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-20 - Cache Key Generation Overhead (Redis)
**Learning:** In Go, dynamically formatting cache keys using `strconv.Itoa` and allocating strings via `hex.EncodeToString` inside a high-frequency method like `RedisCache.key` creates unnecessary memory allocations and CPU overhead in the hot path.
**Action:** Replace `hex.EncodeToString` with a fixed-size stack-allocated array for hex encoding, and hardcode fast paths for common cache dimensions (e.g., 512, 1024) to avoid runtime integer-to-string conversions entirely.
