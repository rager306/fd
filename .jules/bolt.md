## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2026-07-13 - Zero-allocation string-to-byte conversion and stack-allocated hex encoding
**Learning:** `[]byte(text)` and `hex.EncodeToString` in hot paths cause unnecessary heap allocations due to byte slice copying and string allocation.
**Action:** Use `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte conversion (with a check for empty string) and `hex.Encode` into a stack-allocated buffer (e.g. `var buf [64]byte`) followed by a standard string conversion to minimize allocations.
