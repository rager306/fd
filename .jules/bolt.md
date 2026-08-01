## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-12 - Memory Leak via Substring from hex.EncodeToString
**Learning:** `hex.EncodeToString(h[:])[:12]` causes a memory leak because taking a substring keeps the entire original 64-character backing string array alive in memory as long as the 12-character string is referenced. This also requires allocating useless heap space for 64 characters when only 12 are needed.
**Action:** When truncating hex encoded output, encode only the required prefix into a precise stack-allocated buffer (e.g., `var dst [12]byte`, `hex.Encode(dst[:], h[:6])`), and cast directly to a string: `string(dst[:])`.
