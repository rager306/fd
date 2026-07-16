## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-02 - Redis Cache Key Generation Overhead
**Learning:** In Go, creating a cache key for each request frequently allocates on the heap. `hex.EncodeToString` and string-to-byte casting are allocation heavy in hot paths.
**Action:** Use zero-copy `unsafe.Slice(unsafe.StringData(...))` to cast strings to bytes and stack-allocated `var buf [64]byte` with `hex.Encode` followed by string cast instead of `hex.EncodeToString`. Additionally, use fast-path hardcoded values for common dimensions (512, 1024) to skip `strconv.Itoa` overhead.
