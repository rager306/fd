## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-01-20 - Fast String Hashing via Zero-Copy Cast
**Learning:** Using `[]byte(text)` before hashing a string creates a new allocation and copies the string contents. Additionally, using `hex.EncodeToString(h[:])` allocates a new string.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte casting before hashing, and use `hex.Encode` into a stack-allocated byte array before string conversion to avoid allocations.
