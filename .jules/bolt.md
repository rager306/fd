## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Request ID Generation Overhead
**Learning:** Generating UUID-like strings using `fmt.Sprintf` for fallbacks or multiple `hex.EncodeToString` + string concatenations for UUIDs causes multiple heap allocations.
**Action:** Replace `fmt.Sprintf` with manual hex conversion into a stack-allocated byte array for the fallback. For UUIDs, encode directly into a single 36-byte array and return `string(buf[:])` to minimize allocations and CPU overhead.
