## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - ETag Generation Memory Overhead
**Learning:** Generating ETag strings by concatenating a quoted hex-encoded string (e.g. `"` + hex.EncodeToString() + `"`) in HTTP middleware causes multiple unnecessary heap allocations on every request.
**Action:** When constructing strings that wrap hex-encoded byte data with fixed prefixes/suffixes (like quotes), pre-allocate a fixed-size byte array on the stack (e.g. `var dst [66]byte`), assemble the string manually, and convert it to a string once to drastically reduce heap allocations.
