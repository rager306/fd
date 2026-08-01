## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-08-01 - Avoid string subslicing for hex encoding prefixes
**Learning:** In Go, subslicing the string returned by `hex.EncodeToString(h)[:12]` keeps the entire 64-byte original backing array alive in memory because strings in Go share their underlying byte arrays when sliced. In a high-frequency hot-path like a cache hash generator, this leads to significant memory retention/leaks and excessive heap allocations.
**Action:** Instead of encoding the full array and slicing the resulting string, slice the source byte array to the desired length (e.g. `h[:6]`), encode it directly into a fixed-size stack-allocated buffer (e.g. `var dst [12]byte`, `hex.Encode(dst[:], h[:6])`), and cast it to a string. This ensures only the needed bytes are allocated and processed, avoiding memory leaks and reducing GC pressure.
