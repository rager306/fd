## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-24 - RedisCache Key Generation Optimization
**Learning:** In Go, using `hex.EncodeToString` dynamically allocates a string on the heap, and appending multiple strings causes further allocations. In high-frequency hot paths like Redis key generation for cache lookups, this creates measurable GC pressure.
**Action:** Encode hashes directly into stack-allocated byte arrays using `hex.Encode` and fast-path string conversions to minimize heap allocations.
