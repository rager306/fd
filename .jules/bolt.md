## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-31 - Zero-Allocation Hash String Conversion
**Learning:** Using `hex.EncodeToString(hash[:])` to create a hex string of a hash causes multiple heap allocations. First for the slice, and second for the new string returned by the hex encoding.
**Action:** Use a stack-allocated byte array `var dst [64]byte` and `hex.Encode` with `string(dst[:])` to reduce allocations when converting hashes to strings in hot paths.
