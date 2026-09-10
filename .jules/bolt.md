## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Hash Hex Encoding Optimization
**Learning:** `hex.EncodeToString(hash[:])` dynamically allocates a full 64-byte string on the heap to represent a 32-byte sha256 hash. If only a small substring is needed (e.g. `[:12]`), this causes unnecessary memory bloat and retention because Go strings are immutable and the substring holds onto the full 64-byte backing array.
**Action:** Always slice the underlying hash byte array first (`hash[:6]`), and encode directly into a stack-allocated buffer (e.g., `var dst [12]byte; hex.Encode(dst[:], ...); string(dst[:])`) to eliminate these heap allocations in hot paths.
