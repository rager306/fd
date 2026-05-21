# S01: Packaged ONNX smoke proof — UAT

**Milestone:** M039-aexhf5
**Written:** 2026-05-21T11:22:26.554Z

# UAT — M039 S01

Packaged image `fd-api:onnx1024-m039` was built and smoke-tested.

Acceptance signals:

- backend `onnx`;
- artifact verified true;
- tokenizer verified true;
- runtime library verified true;
- provider `CPUExecutionProvider`;
- embedding length 1024;
- norm close to 1;
- isolated cache namespaces used;
- container stopped;
- port 18000 clean.

