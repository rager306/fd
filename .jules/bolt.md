## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-07-17 - Avoid string-to-byte cast and hex.EncodeToString allocations
**Learning:** In hot paths (like cache key generation), standard []byte(string) casts and hex.EncodeToString both force heap allocations. Using unsafe.StringData to alias string bytes, and hex-encoding directly into a stack-allocated buffer (e.g. var buf [64]byte) before converting to string, drastically reduces these allocations.
**Action:** For extremely hot hashing functions, avoid []byte() casting via unsafe and prefer fixed-size stack buffers for encoding before string conversion.
