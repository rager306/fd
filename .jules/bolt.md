## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2023-11-02 - Optimize Cache Hashing Allocations
**Learning:** `hex.EncodeToString(hash[:])[:12]` causes unnecessary heap allocations by allocating a full 64-char string on the heap before slicing. In Go, string-to-byte slice conversions (`[]byte(value)`) also allocate memory on the heap.
**Action:** Use `unsafe.Slice(unsafe.StringData(value), len(value))` for zero-allocation string-to-byte conversion (only when read-only safe) and encode only the required prefix into a stack-allocated array (e.g. `var dst [12]byte; hex.Encode(dst[:], h[:6])`) to minimize GC overhead in hot paths.
