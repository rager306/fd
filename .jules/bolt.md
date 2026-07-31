## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-27 - Hash Truncation Memory Optimization
**Learning:** `hex.EncodeToString` allocates memory for the full hex string. If only a prefix of that string is used (e.g. `hex.EncodeToString(h)[:12]`), the entire backing array stays alive, causing a memory leak.
**Action:** Avoid `hex.EncodeToString` when only needing a prefix. Instead, allocate a smaller stack array (e.g. `var dst [12]byte`), use `hex.Encode(dst[:], h[:6])`, and return `string(dst[:])`. Use `unsafe.Slice` on `unsafe.StringData(value)` to avoid copying string bytes to a `[]byte` during hashing for zero-allocation hashing.
