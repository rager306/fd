## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-28 - Cache Key Hex Encoding Allocations
**Learning:** In Go, using `hex.EncodeToString(h[:])` inside a frequently called hot path (like Redis cache key generation) allocates memory on the heap every time. When constructing highly repetitive strings, encoding directly into a stack-allocated byte array (e.g. `var dst [64]byte`) and converting to a string reduces heap allocations and garbage collection overhead.
**Action:** Replace `hex.EncodeToString` with stack-allocated byte arrays for frequently computed hex representations, and apply fast paths for commonly known append strings to avoid runtime `strconv` conversions.
