## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-12 - Memory leaks from string sub-slicing
**Learning:** Slicing an allocated string (e.g., `hex.EncodeToString(h[:])[:12]`) does not free the underlying memory; it keeps the entire original backing array alive. For a cache holding many short hashes, this leads to a continuous memory leak.
**Action:** When truncating encoded output, avoid intermediate string allocations entirely. Allocate a stack buffer for the exact size needed, encode the required prefix directly into it (`hex.Encode(dst[:], h[:6])`), and cast it to a string.
