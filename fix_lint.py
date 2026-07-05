import subprocess
import re

def run_lint():
    result = subprocess.run(["go", "run", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2", "run", "./..."], cwd="api", capture_output=True, text=True)
    return result.stdout + result.stderr

print(run_lint())
