## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Hex Encoding Slicing Allocations
**Learning:** Using `hex.EncodeToString(hash[:])[:length]` allocates a full string (e.g. 64 chars) on the heap before slicing. For small substrings, this adds unnecessary allocation overhead in hot paths.
**Action:** Allocate a small array of the target size on the stack (`var dst [12]byte`), encode only the required prefix bytes (`hex.Encode(dst[:], h[:6])`), and return it as a string to eliminate the temporary heap allocation.
