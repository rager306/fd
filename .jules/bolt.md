## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2026-10-07 - Branch hoisting in loops
**Learning:** Checking a constant configuration variable (like `encodingFormat`) inside a tight loop processing batch elements adds unnecessary string comparison and branch evaluation overhead.
**Action:** Hoist loop-invariant checks (e.g., `isBase64 := encodingFormat == ...`) outside the loop to eliminate redundant comparisons per element.
