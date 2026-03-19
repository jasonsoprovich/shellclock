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

### Project Tags
- [x] **Tags data model** — `Project` struct has `Tags []string` field; persisted in JSON as `"tags": [...]`, omitted when empty; `UpdateProjectTags()` store helper
- [x] **Tag entry in create/rename flow** — after confirming a project name (new project `N` or rename `E`), a second inline prompt appears for comma-separated tags; `enter` saves (empty clears tags), `esc` cancels
- [x] **Direct tag edit** (`#` on project) — opens the tag edit prompt pre-filled with existing tags; works any time a project is focused in the tree view
- [x] **Tag pills in tree view** — tags render as small teal pill badges (colored background, bold text) next to the project name; selected rows show tags as `[tag1, tag2]` inside the highlight
- [x] **Tag filter in report view** (`f`) — opens a centered modal tag picker listing all unique tags; `↑`/`↓` navigate, `enter` applies the selected tag as a filter, `(show all)` clears it; active filter shown with `✓` marker; `esc`/`f` closes without changing the filter
- [x] **Filter indicator in report header** — when a tag filter is active, a pill + `[f] change or clear` line appears below the separator
- [x] **Filtered totals** — grand total in report header reflects only the filtered projects

### Import

- [x] **Toggl CSV import** (`shellclock import toggl /path/to/export.csv`) — CLI subcommand (no TUI); reads a Toggl CSV export and merges it into the shellclock data file. Column mapping: Toggl **Project** → shellclock Project; Toggl **Task** (falling back to **Description**) → shellclock Task; each row → Session with correct start/end timestamps and duration. Projects and tasks matched by name — existing ones are reused and sessions appended, never duplicated. If Description differs from the task name it is stored as session notes (truncated to 120 chars). Prints `X projects imported, Y tasks imported, Z sessions imported` on success.

### Backup
- [x] **Automatic daily backup** — on every launch, before data is loaded, `model.RunBackup()` copies `shellclock.json` to `~/.config/shellclock/backups/shellclock-YYYY-MM-DD.json`; one backup per calendar day; at most 7 backups retained (oldest pruned automatically); completely silent with no UI interaction required
- [x] **Backup info overlay** (`B` in tree view) — centered modal over dimmed background listing all available backup files newest-first, and showing the full backup directory path; dismissed by `esc` or any key

### Idle Timer Warning

- [x] **Idle warning indicator** — when a timer has been running (not paused) past the configured threshold, a flashing `⚠ Nh+` badge appears next to the timer badge in both the tree view and task detail view. The indicator alternates on/off every second using the existing tick chain. No automatic action is taken — it is purely visual.
- [x] **Configurable threshold** (`W` in tree view) — opens an inline prompt pre-filled with the current threshold in minutes; entering `0` disables the feature. Default: 120 minutes (2 hours). Setting is persisted in the data file under `idle_warn`.
- [x] **`DefaultIdleWarnThresholdMins = 120`** — named constant in `model/model.go`; change it to adjust the default for new installs without touching existing data files.
- [x] **On by default** — new and existing data files without `idle_warn` in JSON initialize to `{enabled: true, threshold_mins: 120}`.

### Summary View

- [x] **Summary view** (`s` from tree view) — shows all sessions logged today (default) or this week, grouped by project → task. Each session row displays its time range (e.g., `10:05 AM — 11:30 AM`) and duration, right-aligned. Project and task subtotals shown. Grand total at the bottom. Toggle between today and this week with `w`. Scroll with `↑`/`↓` or `j`/`k`; dismiss with `esc` or `q`. Empty-state message when no sessions exist for the selected period.
- [x] **Timer key change in tree view** — to free `s` for the summary, the tree-view timer start/pause/resume moved from `s` to `p`. Task detail view is unchanged (`s`/`enter` still controls the timer there).

### Help Screen

- [x] **In-app help screen** (`H` from any view) — full-screen scrollable reference panel covering all keybindings organized by view (Tree, Task Detail, Session Editor, Report, Theme Picker, Summary), backup info, CLI commands, and the Toggl importer. Scroll with `↑`/`↓` or `j`/`k`; dismiss with `esc` or `q`. Accessible globally — `H` is intercepted in the root App model before any view-specific handling.

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
