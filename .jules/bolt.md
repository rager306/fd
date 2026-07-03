## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-24 - Zero-Allocation Hashing
**Learning:** Using `[]byte(string)` and `hex.EncodeToString` in hot paths like cache key generation causes unnecessary heap allocations. `[]byte(string)` allocates a new slice, and `hex.EncodeToString` allocates another string.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte casting, and `hex.Encode` with a fixed-size stack-allocated buffer (e.g. `var buf [64]byte`) followed by `string(buf[:])` to avoid `hex.EncodeToString` allocations.
