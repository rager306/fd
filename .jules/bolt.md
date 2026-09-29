## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Hex Encoding String Slicing Memory Bloat
**Learning:** In Go, string slicing (like `hex.EncodeToString()[:12]`) retains the entire backing array in memory. When generating many short hashes, this causes memory bloat.
**Action:** Use a correctly sized stack-allocated array (e.g., `var dst [12]byte`) and `hex.Encode` directly into it, then convert to a string to avoid retaining the full backing array.
