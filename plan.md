1. **Optimize `shortHash` in `api/cache/hash.go` to reduce memory allocations:**
   The `shortHash` function uses `hex.EncodeToString(h[:])[:12]`, which allocates a 64-character string and then slices it down to 12. We can optimize it by using `hex.Encode` into a fixed 12-byte array using only the first 6 bytes of the hash.

2. **Verify changes and run tests:**
   Run tests in `api/cache` and ensure benchmark performance improvements.

3. **Complete pre-commit steps to ensure proper testing, verification, review, and reflection are done.**

4. **Submit the PR**
   Use the submit tool to push the changes with the required Bolt PR format.
