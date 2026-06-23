## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Computation Allocations
**Learning:** In Go, calling `sha256.Sum256([]byte(value))` causes a heap allocation for the byte slice when `value` is a string. Furthermore, `hex.EncodeToString(h[:])` causes allocations. Using `unsafe.Slice` to cast strings to bytes without allocations, and `hex.Encode` into a stack-allocated array (e.g. `var buf [64]byte`) followed by a standard `string(buf[:])` conversion avoids heap allocations during hashing operations in hot paths.
**Action:** When hashing strings in a hot path, use `unsafe.Slice` for zero-copy string-to-bytes conversion (adding the gosec bypass comment) and use stack-allocated arrays for hex encoding before converting back to strings.
