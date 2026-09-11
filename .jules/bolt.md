## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Encoding Allocation Overhead
**Learning:** In Go, calling `hex.EncodeToString(hash[:])` dynamically allocates a full-length string on the heap, even if the result is immediately sliced (e.g., `[:12]`). Additionally, converting strings to byte slices for hashing forces heap allocations.
**Action:** For frequently called hashing functions (like short IDs or cache keys), perform zero-allocation string conversions using `unsafe.StringData`, and encode hex values directly into a small, stack-allocated byte array (e.g., `var dst [12]byte; hex.Encode(dst[:], hash[:6])`) before converting to a string.
