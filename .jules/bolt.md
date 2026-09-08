## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-01-01 - Avoid hex.EncodeToString and string concatenation for hashes
**Learning:** Generating string representations of hashes using `hex.EncodeToString(hash[:])` and concatenating static characters (`+`) causes multiple dynamic heap allocations.
**Action:** Allocate a fixed-size byte array on the stack (e.g., `var buf [66]byte`), use `hex.Encode` directly into the array, and convert it to a string once (`string(buf[:])`) to minimize allocations and improve speed.
