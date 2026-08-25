## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - ETag Generation Overhead
**Learning:** Generating strings that concatenate dynamic data (like hex encoding) with static literals (like quotes) using standard string concatenation `+` causes multiple heap allocations per operation in Go, which adds significant GC pressure in high-throughput hot paths like middleware.
**Action:** Pre-allocate a byte array of the exact size on the stack, build the string components into it directly (using functions like `hex.Encode` instead of `hex.EncodeToString`), and perform a single `string()` cast at the end to eliminate intermediate allocations.
