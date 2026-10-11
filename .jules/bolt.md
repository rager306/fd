## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.
## 2025-02-14 - Dead Code Elimination in Benchmarks
**Learning:** When benchmarking string operations in Go, assigning the result to the blank identifier (`_ = func()`) allows the compiler's dead-code elimination (DCE) to remove the allocation entirely, resulting in falsely optimistic benchmarks (e.g., 0 allocs/op).
**Action:** Always assign benchmark results to a global, exported variable (e.g., `var GlobalSink string`) to ensure the compiler actually executes and allocates the data structure being tested.
