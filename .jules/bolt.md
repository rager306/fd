## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-07-05 - Avoid memory allocations in hashing operations
**Learning:** `hex.EncodeToString(sha256.Sum256([]byte(string)))` incurs significant allocation overhead due to the byte slice conversion and `EncodeToString`.
**Action:** Use string-to-byte pointer casting via `unsafe.Slice(unsafe.StringData(text), len(text))` along with `hex.Encode` into stack-allocated buffers `var buf [x]byte` to reduce memory allocations in Go when hashing and converting to hexadecimal representation in hot paths.
