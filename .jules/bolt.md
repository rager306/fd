## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-10-25 - ETag String Assembly Optimization
**Learning:** Combining hex-encoded data with string literals using `+` concatenation in high-frequency paths (like middleware ETags) forces multiple heap allocations (3 allocs per call).
**Action:** Assemble the parts directly into a stack-allocated byte array (`var dst [66]byte`) and use a single `string(dst[:])` cast to reduce GC pressure and improve throughput.
