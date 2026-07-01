## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-24 - Zero-Allocation Cache Key Hashing
**Learning:** Hashing string keys in hot paths using `[]byte(text)` and `hex.EncodeToString` creates significant GC overhead due to repeated heap allocations.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy byte casting and a stack-allocated buffer (e.g., `var buf [64]byte`) with `hex.Encode` to avoid unnecessary allocations.
