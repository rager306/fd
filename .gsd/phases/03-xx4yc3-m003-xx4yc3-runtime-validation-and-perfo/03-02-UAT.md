# S02: Live API smoke tests — UAT

**Milestone:** M003-xx4yc3
**Written:** 2026-05-19T08:15:30.697Z

# UAT: S02 Live API smoke tests

## Verification performed

- API health: `api=ok`.
- TEI health: curl HTTP success with empty body.
- Redis ping: `PONG`.
- `/v1/embeddings` 1024d single: `data_len=1`, `dim=1024`, `emb_len=1024`.
- `/v1/embeddings` 512d array: `data_len=2`, `dims=[512,512]`, `emb_lens=[512,512]`.
- `/embeddings/batch` base64 1024d: `count=2`, `dimensions=1024`, valid base64.
- `/embeddings/batch` float 512d: `count=1`, `dimensions=512`, decoded float array length 512.
- Negative cases: invalid JSON, empty input, bad dimensions, bad format all returned 400.

