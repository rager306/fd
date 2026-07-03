# S01 Assessment

**Milestone:** M016-pdcjat
**Slice:** S01
**Completed Slice:** S01
**Verdict:** roadmap-confirmed
**Created:** 2026-05-20T05:36:48.889Z

## Assessment

S01 strongly supports max sequence length 128 as the primary suspect: all 17 worst M015 cases truncate at 128, only 2 at 512. S02 should focus on runtime/vector diagnostics comparing TEI to local ONNX at sequence lengths 128 and 512 on the same IDs.
