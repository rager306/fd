## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-18 - Avoid Memory Bloat When Slicing Hex Strings
**Learning:** Slicing a large hex-encoded string (e.g., `hex.EncodeToString(h[:])[:12]`) keeps the entire backing array (64 bytes for SHA-256) alive in memory, even though only a small prefix is needed. In high-throughput paths, this leads to memory bloat.
**Action:** To extract a short hex string from a hash without retaining the full backing array, allocate a fixed-size stack buffer, use `hex.Encode`, and cast exactly the needed slice to a string: `var buf [64]byte; hex.Encode(buf[:], h[:]); return string(buf[:12])`.
