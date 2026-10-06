## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-10-06 - Redis Key Generation strconv.Itoa Optimization
**Learning:** `strconv.Itoa` causes a heap allocation for integers > 99. In high-throughput paths (like Redis cache key generation for 512d or 1024d embeddings), this adds unnecessary memory overhead. Go's string concatenation (e.g., `a + b + c`) is highly optimized by the compiler into a single allocation when using string literals. Replacing `strconv.Itoa` with fast-path string literals allows for completely optimized concatenations and avoids extra allocs.
**Action:** Use fast-path string literals for known large integer constants when building strings in hot paths, and rely on standard string concatenation instead of manual byte arrays unless necessary.
