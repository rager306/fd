## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-10-31 - Cache Key String Slicing Bloat
**Learning:** In Go, string slicing (e.g., `str[:12]`) retains a reference to the original string's entire backing array. If slicing a large dynamically allocated string (like `hex.EncodeToString()` output), the entire backing array remains in memory, causing memory bloat for long-lived cache keys.
**Action:** For partial hex encoding of hashes, encode directly into a correctly sized stack-allocated array (e.g., `var dst [12]byte; hex.Encode(dst[:], h[:6])`) and convert it to a string. This minimizes allocations and reduces long-term memory footprint.
