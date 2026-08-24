## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - ETag Generation Overhead
**Learning:** Generating string hashes and concatenating them with quotes in high-frequency middleware like `CacheHeaders` causes unnecessary heap allocations. Using `hex.EncodeToString(sum[:])` allocates twice (once for the slice, once for the string).
**Action:** Assemble combined strings directly into stack-allocated byte arrays of exact size (e.g., `var dst [66]byte` for quoted SHA256 hex) and use `hex.Encode` followed by a single cast to `string(dst[:])` to eliminate multi-allocations.
