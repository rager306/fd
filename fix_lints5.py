import re

def fix_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    for search, replace in replacements:
        content = content.replace(search, replace)
    with open(filepath, 'w') as f:
        f.write(content)

fix_file('api/embed/coalescedembedder.go', [
    ('\t\tif err != nil {\n\t\t\tc.logger.Error("coalescedembedder: TEI fallback single query failed", "error", err, "size", len(coalescedTexts))\n\t\t\tfailBatch(itemBatch, err)\n\t\t} else if len(fallbackResp.Embeddings) != 1 {\n\t\t\tc.logger.Error("coalescedembedder: TEI fallback response length mismatch", "expected", 1, "got", len(fallbackResp.Embeddings))\n\t\t\tfailBatch(itemBatch, fmt.Errorf("fallback embeddings length mismatch: %d", len(fallbackResp.Embeddings)))\n\t\t} else {',
    '\t\tswitch {\n\t\tcase err != nil:\n\t\t\tc.logger.Error("coalescedembedder: TEI fallback single query failed", "error", err, "size", len(coalescedTexts))\n\t\t\tfailBatch(itemBatch, err)\n\t\tcase len(fallbackResp.Embeddings) != 1:\n\t\t\tc.logger.Error("coalescedembedder: TEI fallback response length mismatch", "expected", 1, "got", len(fallbackResp.Embeddings))\n\t\t\tfailBatch(itemBatch, fmt.Errorf("fallback embeddings length mismatch: %d", len(fallbackResp.Embeddings)))\n\t\tdefault:')
])

fix_file('api/queue/types.go', [
    ('StatusCompleted Status = "completed"\n\tStatusFailed    Status = "failed"', '// StatusCompleted means job finished successfully.\n\tStatusCompleted Status = "completed"\n\t// StatusFailed means job errored.\n\tStatusFailed    Status = "failed"')
])

fix_file('api/embed/coalesce_baseline_test.go', [
    ('func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (calls, totalTexts int, durations []time.Duration)', 'func runCorpusBurst(t *testing.T, e Embedder, texts []string, concurrency int) (totalTexts int, durations []time.Duration)')
])
