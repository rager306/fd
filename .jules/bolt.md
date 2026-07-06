## 2025-02-12 - Fast Base64 Encoding
**Learning:** `base64.StdEncoding.EncodeToString` allocates a string internally. We can avoid this allocation and speed up encoding by passing a stack-allocated buffer to `base64.StdEncoding.Encode` and returning it as a string without allocation using `unsafe.String(unsafe.SliceData(buf), len(buf))`.
**Action:** Use fixed-size stack-allocated buffers and `base64.StdEncoding.Encode` in hot paths instead of `base64.StdEncoding.EncodeToString`.

## 2025-02-12 - Fast Base64 Encoding
**Learning:** `base64.StdEncoding.EncodeToString` allocates a string internally. We can avoid this allocation and speed up encoding by passing a stack-allocated buffer to `base64.StdEncoding.Encode` and returning it as a string without allocation using `unsafe.String(unsafe.SliceData(buf), len(buf))`.
**Action:** Use fixed-size stack-allocated buffers and `base64.StdEncoding.Encode` in hot paths instead of `base64.StdEncoding.EncodeToString`.
