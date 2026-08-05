## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - hex.EncodeToString Heap Allocation Overhead
**Learning:** In Go, `hex.EncodeToString` allocates a new string and an intermediate byte slice on the heap. In hot paths (like cache hashing or ETag generation), this causes measurable memory allocation overhead.
**Action:** For zero-allocation string building in high-frequency paths, use stack-allocated `[]byte` buffers and `hex.Encode(dst[:], src[:])`, then cast directly to a string.
