---
estimated_steps: 1
estimated_files: 1
skills_used: []
---

# T05: Recommend ONNX CPU optimization benchmark path

Produce ONNX CPU optimization recommendation with provider/quantization/NUMA benchmark matrix, required env/config snapshot fields, success metrics, and stop criteria.

## Inputs

- `S06 T01`
- `S06 T02`
- `S06 T03`
- `S06 T04`

## Expected Output

- `S06 summary and S03 input`

## Verification

Research artifact includes matrix, ranked options, benchmark config fields, and exclusions.

## Observability Impact

Turns ONNX CPU tuning claims into executable future benchmark plan.
