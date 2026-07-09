import re

def fix_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    for search, replace in replacements:
        content = content.replace(search, replace)
    with open(filepath, 'w') as f:
        f.write(content)

fix_file('api/embed/coalesce_baseline_test.go', [
    ('func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (totalTexts int, durations []time.Duration)', '//nolint:unparam // useful for metrics reporting in test\nfunc runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (calls, totalTexts int, durations []time.Duration)')
])
