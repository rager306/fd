## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Zero-allocation String Hashing
**Learning:** In hot paths (like cache lookups), standard string-to-byte slice conversions (`[]byte(str)`) and `hex.EncodeToString` cause multiple heap allocations.
**Action:** Use `unsafe.Slice(unsafe.StringData(str), len(str))` for read-only hashing and stack-allocated `var dst [12]byte` with `hex.Encode` to reduce allocations from 3 to 1 per call.
