## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2024-05-15 - Go Hex Encode String Slicing Bloat
**Learning:** Using `hex.EncodeToString(hash[:])[:12]` causes a heap allocation for the full 64-character hex string, and slicing it retains the entire backing array in memory, causing memory bloat in hot paths.
**Action:** Encode the exact required bytes directly into a correctly-sized stack-allocated array (e.g. `var dst [12]byte; hex.Encode(dst[:], hash[:6])`) to minimize allocations and avoid retaining unneeded backing arrays.
