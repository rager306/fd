## 2024-07-08 - Avoid string allocation when hex encoding
**Learning:** Using `hex.EncodeToString` allocates a new string. When hex encoding in a hot path, encoding directly into a pre-allocated byte array or slice using `hex.Encode` and then converting it to string (or combining with other strings like building a UUID) reduces memory allocations and significantly speeds up the code execution.
**Action:** When a hex-encoded string is needed, use `hex.Encode` into a pre-allocated stack buffer instead of `hex.EncodeToString`.
