## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-07-12 - Hex Encoding Overhead
**Learning:** Using `hex.EncodeToString` requires an extra allocation for the underlying byte slice. Hex-encoding directly into a stack-allocated buffer (e.g. `var buf [64]byte`) and explicitly casting to string reduces allocations and improves throughput by ~10% in hot paths.
**Action:** Replace `hex.EncodeToString` with `hex.Encode` into stack-allocated buffers for hashing operations in hot paths.
