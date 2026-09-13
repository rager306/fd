## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-24 - Hex Encoding Allocations
**Learning:** `hex.EncodeToString` dynamically allocates a string on the heap. String to byte slice conversion `[]byte(string)` also causes heap allocation.
**Action:** For performance critical hashing, use stack-allocated arrays `var dst [12]byte` with `hex.Encode`, and zero-allocation string-to-byte slice conversion `unsafe.Slice(unsafe.StringData(str), len(str))` to avoid allocations.
