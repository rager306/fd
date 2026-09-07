## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-11-20 - Hex Encoding and String to Bytes Allocations
**Learning:** In Go, `hex.EncodeToString(hash[:])[:12]` allocates a full 64-character string on the heap before slicing. Similarly, `[]byte(string)` allocates a new slice. Both are unnecessary in hot paths like cache hashing.
**Action:** For string to byte conversion in read-only operations like hashing, use `unsafe.Slice(unsafe.StringData(str), len(str))` guarded by a length check. To get a short hex hash, encode only the required bytes directly into a stack-allocated byte array (`var dst [12]byte; hex.Encode(dst[:], hash[:6])`) and convert it to string once.
