# M030-3kdha1: ONNX artifact path security remediation

**Vision:** Close the remaining low-severity ONNX operational security findings so the next gates can focus on immutable artifact source selection and hosted workflow proof.

## Success Criteria

- M028 LOW-3 and LOW-4 remediated with code/probes.
- Existing local ONNX artifact paths remain supported.
- No secrets/raw URLs/raw legal text are printed in new artifacts.
- TEI default remains unchanged.
- Milestone committed locally, no push.

## Slices

- [x] **S01: S01** `risk:medium` `depends:[]`
  > After this: After this, manifest/provisioning/verifier artifact paths are root-constrained and diagnostics avoid leaking absolute host paths by default.

- [x] **S02: S02** `risk:medium` `depends:[]`
  > After this: After this, M030 outcome records M028 LOW remediation and final checks pass locally.

## Boundary Map

Not provided.
