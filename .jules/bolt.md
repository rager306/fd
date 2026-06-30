## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-07-01 - Avoid `hex.EncodeToString` in Hot Paths
**Learning:** `hex.EncodeToString` allocates multiple times because it allocates a byte slice for the hex encoded string, and then returns a string (which is a separate allocation).
**Action:** Use `hex.Encode` into a stack-allocated byte array buffer (e.g. `var buf [64]byte`) and then convert that directly to a string. This eliminates intermediate slice allocations and saves memory allocations in hot paths like cache hashing and request ID generation.
