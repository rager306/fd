import re

def fix_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    for search, replace in replacements:
        content = content.replace(search, replace)
    with open(filepath, 'w') as f:
        f.write(content)

fix_file('api/fd_v2_queue_integration_test.go', [
    ('func setupQueueTestServer(t *testing.T, queueCap, batchSize int)', '//nolint:unparam // keep batchSize for parity with production code\nfunc setupQueueTestServer(t *testing.T, queueCap, batchSize int)')
])
fix_file('api/embed/coalesce_baseline_test.go', [
    ('\n\tcalls = 0\n', ''),
    ('calls++\n', '')
])
