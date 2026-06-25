## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-11-06 - Hex Encoding and String Allocation Overhead
**Learning:** In Go, string-to-byte casting (`[]byte(str)`) and encoding entire strings when only a prefix is needed (e.g. `hex.EncodeToString(h[:])[:12]`) create unnecessary memory allocations, which can become bottlenecks in hot paths like cache hashing.
**Action:** Use `unsafe.Slice(unsafe.StringData(str), len(str))` for zero-allocation read-only string-to-byte conversions, and encode directly into stack-allocated arrays (e.g. `var buf [12]byte`) instead of relying on `EncodeToString` when slicing the output.
