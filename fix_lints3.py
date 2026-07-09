with open('api/main.go', 'r') as f:
    content = f.read()

content = content.replace('recoveryCancel()\n\t\t//nolint:gocritic // exitAfterDefer: immediately shutting down on critical cache failure\n\t\tos.Exit(1)', 'os.Exit(1)')

with open('api/main.go', 'w') as f:
    f.write(content)
