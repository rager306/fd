## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-18 - Zero-Allocation Hashing Optimization
**Learning:** Generating hex-encoded hashes using `sha256.Sum256([]byte(text))` and `hex.EncodeToString(h[:])` causes 3 heap allocations (1 for `[]byte(text)`, 1 for slice conversion, 1 for string conversion). This overhead is unnecessary for high-frequency operations like generating cache keys in L2 cache.
**Action:** Use `unsafe.Slice` for string-to-byte-slice conversion when hashing, and stack-allocated byte arrays with `hex.Encode` for hex string generation to reduce allocations down to just 1 for the final returned string.
