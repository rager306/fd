## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Hex Encoding Allocations in Hot Paths
**Learning:** `hex.EncodeToString` allocates both an intermediate byte slice and the final string on the heap. In extremely high-frequency code paths like hash generation for cache keys or ETags, these allocations can add up.
**Action:** Replace `hex.EncodeToString` with `hex.Encode` to a stack-allocated buffer (e.g. `var dst [64]byte`), and cast to string directly. For strings that append literals like quotes (ETags), assemble the entire string (including literals) inside a single stack array first to eliminate intermediate string concatenation allocations entirely.
