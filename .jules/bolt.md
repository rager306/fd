## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2025-02-12 - Reassigning append vs Ignoring append
**Learning:** `append()` in Go returns a potentially new slice header if it reallocates. Modifying an intentional `slice = append(slice, ...)` to `_ = append(slice, ...)` just to suppress a `SA4010` (unused append result) linter warning is extremely dangerous, as it drops the new slice header and prevents the original slice from growing properly, leading to silent data loss or bugs.
**Action:** Always retain the assignment `slice = append(slice, ...)` and suppress the linter warning explicitly with `//nolint:staticcheck // intentional` if the resulting slice is indeed intended to be unused locally but needs to be mutated in scope.
