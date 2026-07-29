## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Zero-Allocation Hex Encoding
**Learning:** In Go, `hex.EncodeToString` allocates a new string and intermediate byte slice on the heap. Using a stack-allocated buffer (`var dst [64]byte`) and `hex.Encode(dst[:], h[:])`, then casting to a string avoids heap allocations.
**Action:** Use stack-allocated buffers and `hex.Encode` for frequent hashing paths (like ETags or Cache Keys) to reduce GC pressure and improve latency.
