import re

def fix_file(filepath, replacements):
    with open(filepath, 'r') as f:
        content = f.read()
    for search, replace in replacements:
        content = content.replace(search, replace)
    with open(filepath, 'w') as f:
        f.write(content)

fix_file('api/embed/coalescedembedder.go', [
    ('if err != nil {\n\t\t\t\tj.result <- coalescedResult{err: err}\n\t\t\t} else if cursor+n <= len(embs) {\n\t\t\t\tslice := make([][]float32, n)\n\t\t\t\tcopy(slice, embs[cursor:cursor+n])\n\t\t\t\tj.result <- coalescedResult{embeddings: slice}\n\t\t\t} else {\n\t\t\t\tj.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}\n\t\t\t}', 'switch {\n\t\t\tcase err != nil:\n\t\t\t\tj.result <- coalescedResult{err: err}\n\t\t\tcase cursor+n <= len(embs):\n\t\t\t\tslice := make([][]float32, n)\n\t\t\t\tcopy(slice, embs[cursor:cursor+n])\n\t\t\t\tj.result <- coalescedResult{embeddings: slice}\n\t\t\tdefault:\n\t\t\t\tj.result <- coalescedResult{err: fmt.Errorf("coalesce split: cursor %d+%d > len(embs) %d", cursor, n, len(embs))}\n\t\t\t}')
])
