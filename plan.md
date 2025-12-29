# zap — project plan

A fast, safe, context-based CLI + TUI for finding and deleting folders (e.g. node_modules, build) using Go.

---

## 1. Project goals (north star)

**zap should be:**

- Damn fast (no unnecessary work)
- Safe by default (context-based, preview before delete)
- Simple UX (one happy path, no config files)
- Keyboard-driven (vim + normie friendly)

**Non-goals (for v1):**

- No global filesystem deletes
- No background indexing or caching
- No config files
- No mouse support
- No delete-from-preview

---

## 2. CLI contract

### Usage

```zsh
zap <folder-name>
zap -r <regex>
zap
```

### Behavior

- `zap <folder-name>`: search immediately
- `zap -r <regex>`: regex-based search (advanced)
- `zap`: open text input to type folder name, then search

### Scope & safety

- Search root = current working directory (cwd)
- No ability to specify `/` or other arbitrary roots
- All paths displayed relative to cwd

---

## 3. High-level UX flow

1. Determine search target (arg, regex, or text input)
2. Start filesystem scan immediately
3. Show spinner **only if scan > 300ms**
4. If no results → print message & exit
5. Show interactive list of matched folders
6. Optional: preview folder contents
7. User selects folders
8. User presses Enter
9. Exit TUI
10. Delete selected folders in parallel
11. Print summary & exit

---

## 4. TUI modes

### Mode 1: List Mode (default)

- Shows search results
- Multi-select
- Entry point for preview
- Only place deletion is allowed

### Mode 2: Preview Mode

- Read-only file tree
- Triggered per-folder
- Used to decide whether to select/delete

---

## 5. Keybindings

### Global navigation (all modes)

- Up: `↑`, `k`, `Ctrl+p`
- Down: `↓`, `j`, `Ctrl+n`

---

### List Mode

**Navigation**

- Top: `gg`, `Home`
- Bottom: `G`, `End`

**Selection**

- Toggle: `Space`
- Select all: `a`
- Deselect all: `A`
- Invert selection: `i`

**Actions**

- Preview folder: `v`, `l`, `Tab`
- Delete selected: `Enter`
- Quit (no delete): `q`, `Esc`, `Ctrl+c`

---

### Preview Mode

**Tree navigation**

- Expand: `Enter`, `l`
- Collapse / back: `h`

**Folder switching**

- Next folder: `n`
- Previous folder: `p`

**Exit**

- Back to list: `q`, `Esc`, `h`

---

## 6. Bottom hint bar (no help overlay)

### List Mode

```zsh
↑↓/jk move • space select • a all • v view • enter delete • q quit
```

### Preview Mode

```zsh
↑↓/jk move • enter/l expand • h back • n/p next • q quit
```

This is the only in-app help for v1.

---

## 7. Filesystem scanning (performance-critical)

### Core rules

- Use `filepath.WalkDir`
- Never call `os.Stat` unnecessarily
- Compare `d.Name()` only (no string contains)
- Prune aggressively with `filepath.SkipDir`

### Matching logic

**Exact match mode:**

- Match when `d.Name() == target`
- Append path
- `SkipDir`

**Regex mode:**

- Compile regex once
- Match against `d.Name()` only
- Append path
- `SkipDir`

### Ignored directories (always pruned)

- `.git`
- `.idea`
- `.vscode`
- `.DS_Store`

---

## 8. Spinner & loading rules

- Start scan immediately in goroutine
- Record start time
- If scan finishes < 300ms → no spinner
- If > 300ms → show spinner
- No fake delays

---

## 9. Preview mode (file tree)

### Trigger

- Only built when user presses preview key
- Never precomputed

### Implementation

- Use `bubbles/filetree`
- Root = selected folder

### Hard limits (for speed & safety)

- Max depth: 3
- Max total items: 300
- If exceeded → show truncated message

### Restrictions

- No deletion in preview mode
- No file sizes
- No gitignore parsing

---

## 10. Deletion phase

### Rules

- Happens only after exiting TUI
- Uses `os.RemoveAll`
- Parallel deletion using goroutines
- Ignore errors silently (best-effort)

### Summary output (after deletion)

- Number of folders deleted
- Total time taken (optional)

---

## 11. Internal architecture

```text
zap/
  cmd/
    root.go        # cobra CLI
  scan/
    scan.go        # filesystem traversal
  tui/
    model.go       # state
    update.go      # key handling
    view.go        # rendering
```

### Separation rules

- `scan` has no TUI imports
- `tui` has no filesystem traversal logic
- `cmd` wires everything together

---

## 12. Build order (recommended)

1. Cobra command (`zap`, args, flags)
2. Scanner package (exact match only)
3. Basic Bubble Tea list view
4. Selection logic
5. Bottom hint bar
6. Deletion logic (post-TUI)
7. Spinner with 300ms cutoff
8. Preview mode (file tree)
9. Regex support (`-r`)
10. Polish & edge cases

---

## 13. Edge cases to handle

- No matches found
- Permission denied (silently skip)
- Very large directories (truncate preview)
- Single match (still show TUI)

---

## 14. MVP exit criteria

zap v1 is complete when:

- Can find and delete folders fast
- No accidental deletes outside cwd
- Preview prevents mistakes
- UX feels instant and predictable
- Binary is small and dependency-light

---

## 15. Future ideas (not now)

- Sort by size
- Size estimation before delete
- Age-based filters
- Configurable ignore list
- Dry-run mode

---

**Guiding principle:**

> zap should feel like `rm`, but with judgment.
