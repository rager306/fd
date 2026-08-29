## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Hash encoding allocations
**Learning:** Using `hex.EncodeToString(hash[:])[:12]` allocates a full 64-character string on the heap before slicing. We can avoid this allocation.
**Action:** Use a stack-allocated byte array (`var dst [12]byte; hex.Encode(dst[:], hash[:6])`) and zero-copy string-to-byte conversion (`unsafe.Slice(unsafe.StringData(str), len(str))`) to minimize overhead and allocations for hot-path hashing.
