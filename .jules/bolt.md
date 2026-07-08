## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-07-08 - Base64 Encoding Optimization
**Learning:** Using `base64.StdEncoding.EncodeToString` creates a string allocation that can be avoided. Additionally, `Float32SliceToBytes` can be optimized heavily on little-endian machines by casting the `[]float32` slice directly to `[]byte` with `unsafe.Slice` and using `copy()`.
**Action:** Use pre-allocated buffers and `base64.StdEncoding.Encode` with `unsafe.String` for zero-allocation base64 string conversion in hot paths. Use `unsafe.Slice` with `copy()` for fast-path endian conversions when dynamically confirmed safe.
