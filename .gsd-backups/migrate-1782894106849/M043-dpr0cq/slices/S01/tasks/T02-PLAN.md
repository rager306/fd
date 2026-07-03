---
estimated_steps: 1
estimated_files: 8
skills_used: []
---

# T02: Fixed all 11 genuine issues: errorlint redis.Nil, gosec G112/G115/G304, revive package-comments + var-naming + early-return, unused consts, test fixture consts

Fix genuine issues found в T01: (a) gosec false positives — добавить //nolint:gosec с justification для legit G107 (URL from env) или fix actual issue. (b) bodyclose — добавить defer resp.Body.Close() если missing. (c) prealloc — preallocate slices с cap. (d) errorlint — заменить errors.New(fmt.Errorf(...)) на fmt.Errorf(...). (e) revive — добавить godoc на exported funcs (WriteError, HTTPStatusFor, AllErrorCodes, EncodeEmbedding, Float32SliceToBytes, BytesToFloat32Slice). Также: убрать underscore в имени parameter если revive flags (`_ = strconv.Itoa` style). Не делать refactor to fix every issue — для genuinely difficult issues (e.g., exported method without doc) добавить //nolint:revive с explicit comment.

## Inputs

- None specified.

## Expected Output

- `api/handlers/embeddings.go (modified)`
- `api/handlers/batch.go (modified)`
- `api/handlers/errors.go (modified)`
- `api/handlers/recovery.go (modified)`
- `api/handlers/notfound.go (modified)`
- `api/middleware/validation.go (modified)`
- `api/embed/codec.go (modified)`
- `api/main.go (modified)`

## Verification

golangci-lint run exit 0. Каждое fix с явным commit reference в fix list. Никаких //nolint:gosec без комментария justification.
