## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-08-10 - Zero-Allocation ETag Generation
**Learning:** String concatenation with `hex.EncodeToString` in high-frequency caching middleware causes excessive heap allocations (3 allocs/op).
**Action:** Use a stack-allocated byte array and `hex.Encode` for zero-allocation string building in hot paths.
