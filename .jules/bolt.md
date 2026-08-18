## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-08-18 - ETag Assembly Optimization
**Learning:** In hot paths (like HTTP middleware for ETags), combining literal strings with hex encoding via standard concatenation (`"` + hex.EncodeToString(...) + `"`) causes multiple allocations.
**Action:** Allocate a fixed-size byte array on the stack, encode directly into the slice `hex.Encode(...)`, set the literal bytes, and perform a single cast to string to save allocations (from 3 down to 1) and improve speed.
