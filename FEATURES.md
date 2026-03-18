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
- [x] **Collapsible project/task tree** — expand/collapse with arrow keys or `h`/`l`
- [x] **Create project** (`N`) — inline text input prompt
- [x] **Create task** (`n`) — inline text input prompt, added under focused project
- [x] **Delete project or task** (`d`) — removes entry and all child data; clamps cursor
- [x] **Time totals** — project and task totals shown inline in the tree
- [x] **Active timer indicator** — `●` marker on the currently timed task
- [x] **Vim-style navigation** — `j`/`k` up/down, `h`/`l` collapse/expand

### Timer View
- [x] **Live HH:MM:SS counter** — tick-based update every second
- [x] **Pause / resume** (`p`) — accumulated time banked on pause; resumes from where it left off
- [x] **Stop & save** (`S`) — commits session to task and clears active timer
- [x] **Reset** (`r`) — discards accumulated time and restarts from now
- [x] **Breadcrumb** — shows `Project › Task` for the active timer
- [x] **Previous total** — shows task's historical total below the clock
- [x] **Paused state styling** — clock and badge switch to yellow when paused
- [x] **Background timer** — returning to tree with `q`/`esc` leaves timer running

### Report View
- [x] **Project and task rows** — hierarchical summary with totals
- [x] **Progress bars** — proportional `█░` bars; project bars relative to grand total, task bars relative to parent project
- [x] **Grand total** — shown right-aligned in the header
- [x] **Active timer notice** — running timer's elapsed time shown at top of report (not yet committed to totals)
- [x] **Scrollable list** — scroll hints appear when content overflows
- [x] **Correct column alignment at any width** — bar/duration columns stay aligned on window resize

### Session Editor
- [x] **Session list** — index, start, end, duration columns for all sessions on a task
- [x] **Add session** (`n`) — form pre-filled with current time; user adjusts with arrow keys
- [x] **Edit session** (`e`) — form pre-filled with existing start/end times
- [x] **Delete session** (`d`) — removes selected session and recalculates totals
- [x] **Second-level precision** — times use format `YYYY-MM-DD HH:MM:SS`
- [x] **Tab / Shift+Tab navigation** — move between start and end fields
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

### UI / Polish
- [x] **Rounded-border panels** — all views use Lip Gloss `RoundedBorder`
- [x] **Context-sensitive help bar** — bottom bar shows keys relevant to the current view and mode
- [x] **Full help toggle** (`?`) — expands help bar to show all bindings
- [x] **Catppuccin Mocha default** — colour palette initialised to Mocha at startup
- [x] **Highlight bleed fix** — selected rows use `Width(w).MaxWidth(w)` to pin highlight to exactly one line regardless of terminal width

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
