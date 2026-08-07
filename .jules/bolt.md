## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-18 - String Sub-slicing Memory Leak
**Learning:** Taking a sub-slice of a string (e.g. hex.EncodeToString(h)[:12]) keeps the entire original backing array alive in memory, causing a memory leak when truncating encoded output.
**Action:** Encode only the required prefix into a precisely-sized stack-allocated buffer (e.g., var dst [12]byte, hex.Encode(dst[:], h[:6])) and cast it directly to a string.
