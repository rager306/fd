## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-24 - Hex Encoding Allocations
**Learning:** Using `hex.EncodeToString(hash[:])` in hot paths causes heap allocations due to dynamic string creation.
**Action:** Encode directly into a stack-allocated byte array (`var dst [64]byte; hex.Encode(dst[:], hash[:])`) and convert to a string (`string(dst[:])`) to minimize allocations in frequent cache key generation.
