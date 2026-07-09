with open('api/main.go', 'r') as f:
    content = f.read()

content = content.replace('os.Exit(1)', 'os.Exit(1) //nolint:gocritic // exitAfterDefer: immediately shutting down on critical startup failure')

with open('api/main.go', 'w') as f:
    f.write(content)
