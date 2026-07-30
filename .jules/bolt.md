## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-11-15 - hex.EncodeToString String Slice Memory Leak
**Learning:** Slicing a string returned by `hex.EncodeToString(h[:])[:12]` prevents the underlying 64-byte backing array from being garbage collected. Additionally, `hex.EncodeToString` performs heap allocations. In high-frequency hot paths like cache hashing, this causes severe memory leaks and excessive allocations.
**Action:** Use a stack-allocated byte array (e.g. `var dst [12]byte`), encode only the required prefix via `hex.Encode(dst[:], h[:6])`, and cast directly to a string: `return string(dst[:])`.
