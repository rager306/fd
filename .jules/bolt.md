## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-28 - Optimize hex.EncodeToString to avoid heap allocations
**Learning:** `hex.EncodeToString(sum[:])` allocates a byte array and string on the heap. String manipulation like string concatenation and slicing (`[:12]`) on these dynamically allocated strings creates unnecessary objects that stay around in memory.
**Action:** Use a stack-allocated buffer (`var dst [64]byte`) in combination with `hex.Encode` and simple casting (`string(dst[:])`) in high-frequency string/hash-generation routines (like `responseETag` and `shortHash`) to improve latency and reduce GC pressure.
