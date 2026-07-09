with open('api/queue/worker.go', 'r') as f:
    content = f.read()

content = content.replace('\tindexByID := make([]*Item, 0, len(batch))\n', '')

with open('api/queue/worker.go', 'w') as f:
    f.write(content)
