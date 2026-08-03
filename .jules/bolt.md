## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-03-09 - Zero-allocation string concatenation for ETags
**Learning:** In Go, combining `hex.EncodeToString` with string concatenation (`"` + str + `"`) allocates multiple times on the heap.
**Action:** For performance-critical hot-paths like HTTP middleware ETag generation, use a single stack-allocated `[66]byte` array to combine leading/trailing literals with hex-encoded data before casting it to a string. This achieves a significant reduction in allocations and memory footprint.
