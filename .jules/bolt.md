## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2024-05-18 - Zero-Allocation Hash Formatting
**Learning:** `hex.EncodeToString` allocates a new byte slice and string every time it's called. Converting a string to a byte slice using `[]byte(text)` also involves an allocation in hot paths.
**Action:** Used `unsafe.Slice(unsafe.StringData(text), len(text))` for zero-copy string-to-byte-slice conversion and replaced `hex.EncodeToString` with `hex.Encode` into a stack-allocated buffer (e.g., `var buf [64]byte`). This eliminated all heap allocations for SHA-256 hash formatting.
