# shellclock — Feature Log

Tracks completed features, active work, and planned additions.
Status key: ✅ Done · 🚧 In Progress · 📋 Planned

---

## ✅ Completed

### Core
- [x] **Project/task data model** — hierarchical: projects → tasks → sessions; single JSON file persistence
- [x] **Atomic file writes** — data saved via tmp-file + rename to prevent corruption
- [x] **Active timer persistence** — timer survives app restarts via `active_timer` in JSON
- [x] **One-timer constraint** — only one timer can run at a time

### Tree View
- [x] **Lipgloss wordmark logo** — rounded-border wordmark (`◷ shellclock`) rendered with lipgloss using `colorOverlay` border and `colorMauve` bold text; adapts to theme changes via `ApplyTheme`; followed by a one-line summary: current date · project count · task count · total tracked time
- [x] **Collapsible project/task tree** — expand/collapse with arrow keys or `h`/`l`
- [x] **Create project** (`N`) — inline text input prompt
- [x] **Create task** (`n`) — inline text input prompt, added under focused project
- [x] **Delete project or task** (`d`) — removes entry and all child data; clamps cursor
- [x] **Rename project or task** (`E`) — inline text input pre-filled with current name; edit and confirm with enter
- [x] **Time totals** — project and task totals shown inline in the tree
- [x] **Live timer overview** — active task row shows `● 00:05:32` / `⏸ 00:05:32`; status line below title shows project › task + elapsed at all times
- [x] **Timer controls from tree** — `s` start/pause/resume, `S` stop & save, `r` reset without leaving the overview
- [x] **Task detail navigation** — `enter` opens task detail; `e` opens dedicated session editor
- [x] **Vim-style navigation** — `j`/`k` up/down, `h`/`l` collapse/expand
- [x] **App-owned tick chain** — tick propagates correctly across all views; no double-ticking

### Task Detail View
- [x] **Opens on `enter`** — pressing enter on any task in the tree opens its detail view
- [x] **Live timer display** — shows `● RUNNING HH:MM:SS` or `⏸ PAUSED HH:MM:SS` with wall-clock start time
- [x] **Timer controls** — `s`/`enter` start/pause/resume, `S` stop & save, `r` reset
- [x] **Cross-task awareness** — shows which other task has the active timer if this task is not the active one
- [x] **Session list** — full scrollable session list with start, end, duration columns
- [x] **Inline session management** — `n` add, `e` edit, `d` delete sessions without leaving the view
- [x] **Paused state styling** — timer badge and elapsed switch to yellow when paused

### Report View
- [x] **Project and task rows** — hierarchical summary with totals
- [x] **Progress bars** — proportional `█░` bars; project bars relative to grand total, task bars relative to parent project
- [x] **Grand total** — shown right-aligned in the header
- [x] **Active timer notice** — running timer's elapsed time shown at top of report (not yet committed to totals)
- [x] **Session notes in report** — sessions with notes show a dimmed `↳ note` line beneath the task row
- [x] **Scrollable list** — scroll hints appear when content overflows
- [x] **Correct column alignment at any width** — bar/duration columns stay aligned on window resize
- [x] **Report export** (`x`) — inline menu to choose CSV or plain text; file saved to `~/.config/shellclock/reports/shellclock-report-YYYY-MM-DD.csv` / `.txt` (directory auto-created); CSV is session-level with columns: Project, Task, Start, End, Duration (h:mm:ss), Notes; plain text mirrors the report view layout; confirmation message shows the full path after successful save

### Session Editor
- [x] **Session list** — index, start, end, duration columns for all sessions on a task
- [x] **Add session** (`n`) — form pre-filled with current time; user adjusts with arrow keys
- [x] **Edit session** (`e`) — form pre-filled with existing start/end times and note
- [x] **Delete session** (`d`) — removes selected session and recalculates totals
- [x] **Second-level precision** — times use format `YYYY-MM-DD HH:MM:SS`
- [x] **Session notes** — optional note field (max 120 chars) on every session; editable in the add/edit form via a dedicated note input (Tab/Shift+Tab/Enter to navigate); note is persisted in JSON and shown in the report view
- [x] **Tab / Shift+Tab navigation** — cycle through start, end, and note fields
- [x] **Validation** — rejects invalid formats and end-before-start
- [x] **Running total** — shown below the session list
- [x] **Layout stability** — scroll hint and error line always reserve fixed height to prevent jitter

### Theme System
- [x] **`Theme` struct** — 13 semantic colour fields; fully decoupled from hex values
- [x] **`ApplyTheme()`** — updates all global colour vars and rebuilds all named styles at runtime
- [x] **31 built-in themes** — Catppuccin family (Mocha, Macchiato, Frappé, Latte), Ayu, Dracula, Everforest, Flexoki, GitHub, Gruvbox, Kanagawa, Material Palenight, Monokai Pro, Nightfox, Nord, One Dark Pro, One Light, Oxocarbon, Poimandres, Rosé Pine family, Solarized, Tokyo Night
- [x] **Live preview** — theme applies immediately on up/down in picker; Esc reverts
- [x] **Theme persistence** — selected theme name stored in JSON; restored on next launch
- [x] **Colour swatches** — each theme row shows 7 representative colour squares in its own colours

### Backup
- [x] **Automatic daily backup** — on every launch, before data is loaded, `model.RunBackup()` copies `shellclock.json` to `~/.config/shellclock/backups/shellclock-YYYY-MM-DD.json`; one backup per calendar day; at most 7 backups retained (oldest pruned automatically); completely silent with no UI interaction required
- [x] **Backup info overlay** (`B` in tree view) — centered modal over dimmed background listing all available backup files newest-first, and showing the full backup directory path; dismissed by `esc` or any key

### UI / Polish
- [x] **Rounded-border panels** — all views use Lip Gloss `RoundedBorder`
- [x] **Context-sensitive help bar** — bottom bar shows keys relevant to the current view and mode
- [x] **Full help toggle** (`?`) — expands help bar to show all bindings
- [x] **Catppuccin Mocha default** — colour palette initialised to Mocha at startup
- [x] **Highlight bleed fix** — selected rows manually pre-padded to content width so highlight never wraps onto a second line
- [x] **Delete confirmation modal** — pressing `d` on any project, task, or session shows a centered overlay modal with the item name and `[y] confirm / [n] / [esc] cancel`; underlying view is dimmed; nothing deleted until `y` is pressed

---

## 🚧 In Progress

_Nothing currently in progress._



---

## 📋 Planned

_No features queued yet. Add new requests here._

---

## Notes

- Feature requests go under **Planned** when described, move to **In Progress** when work starts, and move to **Completed** when merged and functioning correctly.
- The README.md is kept in sync with completed features.
