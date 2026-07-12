## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-12 - Stack allocation for hex encoding in hot paths
**Learning:** Using `hex.EncodeToString` on slices generates unnecessary string allocations, which adds overhead in hot paths (like request IDs, ETags, and cache hashing). If these strings are then concatenated, it causes even more allocations.
**Action:** Avoid `hex.EncodeToString`. Prefer stack-allocating a fixed-size byte array (e.g., `var buf [64]byte`) and use `hex.Encode(buf[:], data)` followed by a standard string conversion (`string(buf[:])`). Embed required extra characters directly into the buffer before conversion to minimize allocations.
