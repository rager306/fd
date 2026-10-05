## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-24 - Hex Encoding and String Slicing Memory Bloat
**Learning:** In Go, string slicing (e.g., `str[:12]`) retains a reference to the original string's entire backing array. Slicing a large dynamically allocated string like the output of `hex.EncodeToString(h[:])` keeps the entire 64-byte allocation in memory.
**Action:** When a shortened hex string is needed, only encode the required number of bytes (e.g., `hex.EncodeToString(h[:6])`) to avoid memory bloat and reduce allocations.
