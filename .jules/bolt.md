## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-10-27 - Hex encoding overhead
**Learning:** Using `hex.EncodeToString(hash[:])[:12]` dynamically allocates a string on the heap twice (once for full hash string, once for truncation).
**Action:** Encode the required bytes directly into a stack-allocated array `var dst [12]byte; hex.Encode(dst[:], h[:6])` and return `string(dst[:])` to minimize allocations and CPU overhead. Also `unsafe.Slice` can avoid the `[]byte(string)` allocation for read-only operations like hashing.
