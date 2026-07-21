## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-23 - Request ID Generation Overhead
**Learning:** In Go, using `hex.EncodeToString` combined with string concatenations in hot paths (like middleware UUID generation on every request) causes unnecessary heap allocations.
**Action:** Replace `hex.EncodeToString` with a stack-allocated byte array and `hex.Encode`, and construct the final string using a zero-allocation byte slice conversion `string(buf[:])`. Similarly, avoid `fmt.Sprintf` for numerical conversions in hot paths by using `strconv.FormatInt` and string concatenation.
