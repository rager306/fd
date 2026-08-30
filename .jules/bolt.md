## 2023-10-27 - Cache Key Generation Overhead
**Learning:** In Go, using `fmt.Sprintf` for constructing strings in highly-frequent hot paths (like cache lookups per embedding input) causes measurable overhead due to reflection and interface boxing, adding unnecessary allocations compared to standard string concatenation.
**Action:** Replace `fmt.Sprintf` with `strconv.Itoa` and simple string concatenation `+` in hot paths, and consider adding fast-path hardcoded values for frequently used parameters (e.g. dimensions 512, 1024) to avoid string conversion entirely.

## 2023-10-27 - Struct Literals in Tests
**Learning:** When adding a pre-computed optimization field (like a cached string prefix) to a struct, unit tests that bypass constructor functions and initialize the struct directly via struct literals (`&RedisCache{...}`) will fail if they rely on the new field.
**Action:** Always search test files for struct literal instantiations and update them to include the newly required optimization fields to ensure the test suite passes on the first run.
