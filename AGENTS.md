# Code Intelligence — native tools

Use the harness's native tools for code navigation, blast radius, and change analysis: **`graph`** (symbol/call/import graph on disk), **`lsp`** (language server + symbol layer), **`edit`** (file edits).

`gitnexus` 1.6.12 is installed and `fd` is indexed (`gitnexus analyze --index-only`; store `.gitnexus/`, kept out of git by `.git/info/exclude`). Its MCP tools (`impact`, `context`, `query`, `detect_changes`) and CLI (`gitnexus ...`, `node .gitnexus/run.cjs ...`) work and stay available as a second opinion — mainly `impact` for a risk verdict and `detect_changes` before committing.

Two caveats in this repo:

- GitNexus hardcodes a bare-name ignore for directories called `cache` (`dist/config/ignore-service.js`, `DEFAULT_IGNORE_LIST`), which hides the whole `api/cache/` package — 14 Go files, and the target of most of issue #480. The repo root `.gitnexusignore` carries `!cache/` to re-include it; keep that file, and check that a package is indexed before trusting an empty result.
- `codebase-memory-mcp` does not start here (MCP `initialize` times out), so do not plan around it.

Re-index after pulling or rebasing: `gitnexus analyze --index-only` (about 40 s, writes no project files).

## Tool map

| Need | Tool call |
| --- | --- |
| Build / refresh the index | `graph` `index` (full), `reindex` (incremental repair), `status` (store counts) |
| How much of the tree the index has seen | `graph` `coverage` |
| Find a symbol | `graph` `search`; workspace-wide `lsp` `symbols` with a `query` (filter the results to `/root/fd` — the module cache is included); `lsp` `find_symbol` inside one file |
| Callers of a symbol | `graph` `trace` with `direction: "in"` |
| Callees of a symbol | `graph` `trace` with `direction: "out"` |
| Exact references and implementations | `lsp` `references`, `find_referencing_symbols`, `implementations`, `find_implementations` |
| Definition, type, hover, file symbols | `lsp` `definition`, `type_definition`, `hover`, `symbols` |
| Architecture, layers, clusters, cycles, hotspots, entry points | `graph` `architecture` with a `view` (returns a summary; go to `graph` `search` / `trace` for per-symbol detail) |
| Git diff → affected symbols | `graph` `changes` |
| Compile / type diagnostics | `lsp` `diagnostics` (one file, or `*` for the whole workspace) |
| Rename a symbol | `lsp` `rename` (with `apply: true`) |
| Edit a whole symbol | `lsp` `replace_symbol_body`, `insert_before_symbol`, `insert_after_symbol`, `safe_delete_symbol` |
| Edit text | `edit` `replace` / `patch` / `hashline` |
| Second opinion on blast radius / diffs | gitnexus MCP `impact`, `context`, `detect_changes`, or `node .gitnexus/run.cjs impact <symbol>` |

## Always Do

- **MUST know the callers before editing a symbol.** Run `graph` `trace` with `direction: "in"` (depth 2) plus `lsp` `find_referencing_symbols`; state the callers you saw, and the confidence, in the change description or PR.
- **MUST treat an empty or low-confidence result as UNKNOWN, not as "unused".** `trace` resolves unresolved call edges by name, so it reports matches with low confidence (~0.4 in this repo) and reports nothing at all for interface calls, dynamic dispatch, reflection, or cross-language calls. An empty caller set is not evidence the symbol is unused. Confirm with `lsp` `references` and a text search before changing or deleting it.
- **MUST refresh the index before trusting it.** Check `graph` `status`; if counts are zero, or the tree moved since the last build, run `graph` `index` with `force: true` — the non-force path is known to leave a reduced edge set after file edits (`rager306/allspark#2`). Read `graph` `coverage`: `partial`, `unusable`, and `skipped` files are unseen, not unaffected.
- **MUST run change analysis before committing.** `graph` `changes` (optionally with `ref_a` / `ref_b`) maps the diff to affected symbols; `lsp` `diagnostics` with `*` checks the whole workspace; then run the project checks: `cd api && go test ./...`, golangci-lint v2.12.2 with `../.golangci.yml`, `go vet ./...`, `govulncheck ./...`.
- **MUST analyze callers, dependencies, and flows at the symbol level.** Text search is the fallback for empty or unknown results and for literals, and the confirmation step for low-confidence matches — never the substitute for the graph walk.
- **MUST use the symbol-aware tools for symbol-level change**: `lsp` `replace_symbol_body` / `insert_after_symbol` / `safe_delete_symbol`, and `lsp` `rename` for renames.

## Never Do

- NEVER edit a function, method, or type before checking its callers.
- NEVER read an empty or low-confidence trace as an all-clear; that is the one verdict that requires confirming by other means.
- NEVER rename symbols with find-and-replace — use `lsp` `rename` with `apply: true`.
- NEVER commit before `graph` `changes` and `lsp` `diagnostics`.
- NEVER report analysis from a stale index as current: if `graph` `status` / `coverage` shows the tree is ahead of the index, re-index, or say plainly in the commit or PR that the analysis predates the change.

## Reporting defects in the harness's own tools

- Defects in the harness itself — the `graph` / `lsp` / `edit` tools, the Python REPL kernel, subagents, MCP connections, compaction — are **not** `fd` issues. Document them in detail and file them in **`rager306/allspark`**: <https://github.com/rager306/allspark/issues>. That repository holds the harness code; `rager306/fd` issues are for this service only.
- Document enough for someone to reproduce without this session: harness version and store paths, project and commit, the exact tool calls, the exact observed output, what you expected, the impact on the task, what you tried, and any workaround.
- Reproduce before filing: run the call more than once and state how often it reproduces. Include the contrast with a tool that answers the same question correctly, when there is one.
- If a reproduction needs a temporary file edit, revert it and verify with `git hash-object <path>` against `git rev-parse HEAD:<path>` plus an empty `git diff`, then say in the report that the edit was temporary.
- If the defect blocks work here, record the local workaround next to the rule it affects (for example: run `graph index force=true` after edits before trusting `graph trace`).

## Index notes

- `graph` `status` prints the store path and node/edge/file counts (one store per project, outside the repo).
- `graph` `index` re-parses the tree; `reindex` repairs the incremental closure. Run `index` after large moves, or when `status` reports 0 nodes.
- A fresh clone has no index — build it before the first navigation query.
