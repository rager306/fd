## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-10-27 - Zero-Allocation Hash Formatting
**Learning:** In highly-frequent paths generating formatted hashes (like `shortHash`), using standard `[]byte(string)` conversion and `hex.EncodeToString()` causes heap allocations due to escaping strings.
**Action:** Use `unsafe.StringData` and `unsafe.Slice` for zero-allocation byte slice conversion when the slice is only read (guarded by a length check). For encoding, use `hex.Encode` into a stack-allocated byte array (`var dst [12]byte`) and return the string conversion to eliminate dynamic string heap allocations.
