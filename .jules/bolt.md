## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-01-20 - Hex Encoding Overhead
**Learning:** In Go, using `hex.EncodeToString(hash[:])[:12]` causes a large heap allocation for the full encoded string before slicing. Additionally, standard `[]byte(string)` conversion allocates.
**Action:** For performance-critical hashing functions (like `shortHash`), use a fixed-size stack array `var dst [12]byte` to encode only the required bytes, and use `unsafe.StringData` for zero-allocation string-to-byte read-only conversion.
