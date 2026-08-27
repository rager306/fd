## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Generation Overhead
**Learning:** The `shortHash` function in `api/cache/hash.go` caused unnecessary allocations on every cache lookup by casting `[]byte(string)` and generating a 64-byte hex string via `hex.EncodeToString` before slicing.
**Action:** Replaced with zero-allocation `unsafe.Slice(unsafe.StringData(value), len(value))` and encoded directly into a stack-allocated byte array (`var dst [12]byte`), dropping allocations from 3 to 1 per call.
