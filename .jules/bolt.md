## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Zero-Allocation Hash Prefixing
**Learning:** In Go, string slicing (e.g., `hex.EncodeToString(h[:])[:12]`) retains a reference to the entire original backing array, causing memory bloat and unnecessary heap allocations. Converting strings to `[]byte` for hashing also causes a heap allocation.
**Action:** Use `unsafe.Slice(unsafe.StringData(str), len(str))` for zero-allocation read-only string-to-byte conversion in extreme hot paths. Pre-allocate stack arrays (`var dst [12]byte`) and use `hex.Encode` to avoid the `EncodeToString` heap allocation and slicing overhead entirely.
