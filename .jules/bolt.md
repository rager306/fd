## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Hex Encoding Overhead
**Learning:** Using `hex.EncodeToString` dynamically allocates a large string on the heap, causing unnecessary memory allocation and bloat in Go hot paths.
**Action:** Avoid it by hex-encoding into correctly sized stack-allocated fixed byte arrays (e.g., `var dst [64]byte`) and converting to a string to minimize allocations and overhead.
