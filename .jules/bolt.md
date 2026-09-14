## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-01-01 - HashText Zero Allocation
**Learning:** Cache key generation overhead due to dynamic string allocations can be effectively mitigated in read-only paths by converting strings to bytes via unsafe slice conversions `unsafe.Slice(unsafe.StringData(text), len(text))` and avoiding `hex.EncodeToString` by encoding directly into stack-allocated byte arrays.
**Action:** Use unsafe slice conversion cautiously, primarily in well-understood, high-frequency, read-only hotspots, explicitly guarding against zero-length strings.
