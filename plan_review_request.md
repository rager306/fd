1. **Analyze `api/middleware/auth.go` timing attack vulnerability**
   - The `APIKeyAuth` middleware uses `subtle.ConstantTimeCompare([]byte(token), []byte(apiKey))` to check the API key.
   - However, `subtle.ConstantTimeCompare` returns `0` immediately if the lengths of the two byte slices are different.
   - This exposes a timing side-channel that leaks the length of the secret `apiKey`.
2. **Fix the timing attack**
   - Instead of comparing the raw byte slices of potentially different lengths, hash both the `token` and `apiKey` using `crypto/sha256.Sum256`.
   - Then, compare the resulting fixed-length hashes (which are always 32 bytes) using `subtle.ConstantTimeCompare`.
   - This prevents any timing differences based on the input length.
3. **Run tests**
   - Run `go test ./...` in the `api` directory to ensure no tests are broken by the change.
4. **Complete pre-commit steps**
   - Complete pre-commit steps to ensure proper testing, verification, review, and reflection are done.
5. **Submit the fix**
   - Submit the PR with the title: `🛡️ Sentinel: [CRITICAL] Fix API Key Timing Attack`
   - Description matching Sentinel's template.
