# zap

A fast, interactive CLI for finding and deleting folders (e.g. `node_modules`, `build`, `.cache`).

## Installation

```bash
go install github.com/coeeter/zap@latest
```

Or build from source:

```bash
git clone https://github.com/coeeter/zap.git
cd zap
go build -o zap .
```

## Usage

```bash
zap                    # Interactive prompt (default: node_modules)
zap <folder-name>      # Search for exact folder name
zap -s <pattern>       # Search with glob pattern
```

### Examples

```bash
zap node_modules       # Find all node_modules folders
zap dist               # Find all dist folders
zap -s "build*"        # Find folders matching build*
zap                    # Opens prompt, defaults to node_modules
```

## Keybindings

### List Mode

| Key           | Action           |
| ------------- | ---------------- |
| `↑` `k`       | Move up          |
| `↓` `j`       | Move down        |
| `gg` `Home`   | Go to top        |
| `G` `End`     | Go to bottom     |
| `Space`       | Toggle selection |
| `a`           | Select all       |
| `A`           | Deselect all     |
| `i`           | Invert selection |
| `v` `l` `Tab` | Preview folder   |
| `Enter`       | Delete selected  |
| `q` `Esc`     | Quit             |

### Preview Mode

| Key         | Action             |
| ----------- | ------------------ |
| `↑` `k`     | Move up            |
| `↓` `j`     | Move down          |
| `Enter` `l` | Expand folder      |
| `h`         | Collapse / go back |
| `n`         | Next folder        |
| `p`         | Previous folder    |
| `q` `Esc`   | Back to list       |

## Features

- **Fast** — Uses `filepath.WalkDir` with aggressive pruning
- **Safe** — Only searches within current directory, preview before delete
- **Interactive** — Vim-style navigation, multi-select, folder preview
- **Parallel deletion** — Deletes folders concurrently

## How It Works

1. Scans current directory for matching folders
2. Shows interactive list for selection
3. Optional: preview folder contents before deciding
4. Deletes selected folders in parallel
5. Shows summary
