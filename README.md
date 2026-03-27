# shellclock

A terminal-based time tracker with a clean TUI, inspired by Toggl. Organize work into projects and tasks, run a live timer, and review time totals — all from the command line without leaving your terminal.

Built with the [Charm](https://charm.sh) stack: Bubble Tea, Lip Gloss, and Bubbles.

---

## Features

- **Projects & Tasks** — hierarchical tree: projects contain tasks, tasks accumulate sessions
- **Live Timer** — start, pause, resume, reset, and stop a timer on any task; timer persists across restarts
- **Session Editor** — view all recorded sessions for a task; add, edit, or delete individual sessions with second-level precision
- **Summary View** — quick daily, weekly, or monthly breakdown of sessions, grouped by project and task
- **Report View** — visual summary with progress bars and percentage breakdown showing time distribution across all projects and tasks
- **Project Tags** — tag projects with arbitrary labels; filter the report to a specific tag
- **Hourly Rate & Earnings** — set a rate per project; toggle an earnings column in the report
- **Idle Timer Warning** — flashing badge when a timer has been running past a configurable threshold
- **Automatic Backups** — daily backup on every launch; last 7 kept; browse and restore from within the app
- **Theme System** — 31 built-in themes (Catppuccin family, Dracula, Nord, Tokyo Night, Gruvbox, Rosé Pine, and more); live preview in the picker; selection persists across launches
- **In-app Help** — press `H` from anywhere to open a full scrollable reference covering every key, view, and CLI command
- **Toggl Import** — import a Toggl Detailed CSV export from the command line without opening the TUI
- **Project Sorting** — cycle sort order in the tree by name or recent use (`o`)
- **Mouse & Scroll Support** — mouse and trackpad scroll work throughout the TUI
- **Keyboard-driven** — vim-style navigation throughout; context-sensitive help bar always visible
- **No external dependencies** — single JSON file, no database, no network

---

## Installation

### Homebrew (macOS and Linux)

```bash
brew install jasonsoprovich/shellclock/shellclock
```

### Build from source

```bash
git clone https://github.com/jasonsoprovich/shellclock.git
cd shellclock
go build -o shellclock .
```

Move the binary somewhere on your `$PATH`:

```bash
mv shellclock /usr/local/bin/
```

### Run directly

```bash
go run .
```

---

## Data storage

All data is stored in a single JSON file:

| Platform | Path |
|----------|------|
| macOS    | `~/Library/Application Support/shellclock/shellclock.json` |
| Linux    | `~/.config/shellclock/shellclock.json` |

The file is written atomically on every change. The active timer and selected theme are both persisted so they survive restarts.

---

## Usage

Launch shellclock:

```bash
shellclock
```

### Importing from Toggl

shellclock can import a Toggl **Detailed** CSV export directly from the command line — no need to open the TUI.

```bash
shellclock import toggl /path/to/export.csv
```

To get the right file from Toggl:
1. Open Toggl Track → **Reports** → **Detailed**
2. Set your date range
3. Click **Export** → **Download as CSV**

On success, shellclock prints a summary and exits:

```
3 projects imported, 4 tasks imported, 20 sessions imported
```

**How data is mapped:**

| Toggl column | shellclock |
|---|---|
| Project | Project (created if new, reused if name matches) |
| Task | Task name (falls back to Description if blank) |
| Description | Task name (if Task is blank); otherwise stored as session notes |
| Start date + Start time | Session start |
| End date + End time | Session end |

Sessions are always appended — existing projects and tasks are matched by name, never duplicated.

---

### Views

shellclock has several views, each accessible from the main tree. Press `H` from any view to open the full in-app help reference.

---

### Tree View (main screen)

The collapsible project/task list. When a timer is running, the active task and elapsed time are shown below the title so you can monitor progress at a glance without leaving the overview.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `→` / `l` | Expand project |
| `←` / `h` | Collapse project |
| `N` | New project |
| `n` | New task (under focused project) |
| `E` | Rename focused project or task |
| `d` | Delete focused project or task |
| `#` | Edit tags for focused project |
| `enter` | Open task detail view |
| `e` | Open session editor for focused task |
| `p` | Start / pause / resume active timer |
| `S` | Stop active timer and save session |
| `r` | Reset active timer |
| `s` | Open summary view (today's sessions) |
| `R` | Open report view |
| `T` | Open theme picker |
| `$` | Set hourly rate for focused project |
| `B` | Show backup info |
| `o` | Cycle project sort order (name ↑/↓, recent ↓/↑) |
| `W` | Idle warn settings |
| `X` | System reset (requires typing CONFIRM) |
| `H` | Open help screen |
| `q` | Quit |
| `?` | Toggle compact key reference bar |

The active task row in the tree always shows the live elapsed time (`● 00:05:32` or `⏸ 00:05:32` when paused).

---

### Task Detail View

Opens when you press `enter` on a task. Shows the live timer and full session list in one place.

| Key | Action |
|-----|--------|
| `s` / `enter` | Start / pause / resume timer |
| `S` | Stop timer and save session |
| `r` | Reset timer |
| `↑` / `k` | Move up in session list |
| `↓` / `j` | Move down in session list |
| `n` | Add new session |
| `e` | Edit selected session |
| `d` | Delete selected session |
| `q` / `esc` | Return to tree |
| `?` | Toggle full key list |

The timer section shows `● RUNNING HH:MM:SS` or `⏸ PAUSED HH:MM:SS` and the wall-clock start time. If a timer is running on a different task, the view shows which task it belongs to instead.

---

### Summary View

Press `s` from the tree to see sessions logged today (default), this week, or this month (cycle with `w`). Sessions are grouped by project and task with time ranges and durations. A grand total appears at the bottom.

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `w` | Cycle between today, this week, and this month |
| `esc` / `q` | Return to tree |

---

### Report View

Summary of all projects and tasks with time totals and proportional progress bars.

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `f` | Filter by tag (tag-picker overlay) |
| `$` | Toggle earnings column |
| `x` | Export report (CSV or plain text) |
| `R` / `q` / `esc` | Return to tree |
| `?` | Toggle full key list |

If a timer is currently running, a notice showing the active task and elapsed time appears at the top of the report.

---

### Session Editor

Per-task session list. Each session has a start time, end time, and calculated duration.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `n` | Add new session |
| `e` | Edit selected session |
| `d` | Delete selected session |
| `q` / `esc` | Return to tree |
| `?` | Toggle full key list |

**Adding or editing a session:**

The form pre-fills with the current time (for new sessions) or the existing times (for edits). Both fields use the format `YYYY-MM-DD HH:MM:SS`.

| Key | Action |
|-----|--------|
| `←` / `→` | Move cursor within field |
| `tab` | Move to next field |
| `shift+tab` | Move to previous field |
| `enter` | Confirm (advances to next field, or saves on the end field) |
| `esc` | Cancel and discard changes |

---

### Theme Picker

Browse and preview all 31 built-in themes. Themes update live as you navigate.

| Key | Action |
|-----|--------|
| `↑` / `k` | Previous theme (live preview) |
| `↓` / `j` | Next theme (live preview) |
| `enter` | Apply and save |
| `esc` / `q` | Cancel (reverts to previous theme) |

**Available themes:**

| Dark | Light |
|------|-------|
| Catppuccin Mocha | Catppuccin Latte |
| Catppuccin Macchiato | Ayu Light |
| Catppuccin Frappé | Everforest Light |
| Ayu Mirage | Flexoki Light |
| Dracula | GitHub Light |
| Everforest Dark | Gruvbox Light |
| Flexoki Dark | One Light |
| GitHub Dark | Rosé Pine Dawn |
| Gruvbox | Solarized Light |
| Kanagawa Wave | |
| Kanagawa Dragon | |
| Material Palenight | |
| Monokai Pro | |
| Nightfox | |
| Nord | |
| One Dark Pro | |
| Oxocarbon | |
| Poimandres | |
| Rosé Pine | |
| Rosé Pine Moon | |
| Solarized Dark | |
| Tokyo Night | |

---

### Help Screen

Press `H` from **any view** to open the full in-app help reference. It covers every keybinding organized by view, backup details, CLI commands, and the Toggl import workflow.

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
| `esc` / `q` | Close and return to tree |

---

## Project structure

```
shellclock/
├── main.go                      # Entry point; routes import subcommand or launches TUI
├── internal/
│   ├── model/
│   │   ├── model.go             # Data types and JSON persistence
│   │   └── backup.go            # Automatic daily backup logic
│   ├── importer/
│   │   └── toggl.go             # Toggl CSV import logic
│   ├── ui/
│   │   ├── app.go               # Root Bubble Tea model, view routing, tick chain
│   │   ├── tree.go              # Project/task tree view (live timer overview)
│   │   ├── taskdetail.go        # Task detail view (timer + session management)
│   │   ├── summary.go           # Daily/weekly session summary view
│   │   ├── report.go            # Summary report view with earnings
│   │   ├── edit.go              # Dedicated session editor view
│   │   ├── export.go            # Report export (CSV / plain text)
│   │   ├── themepicker.go       # Theme selection view
│   │   ├── help.go              # Full-screen help reference view
│   │   ├── confirm.go           # Delete confirmation modal overlay
│   │   ├── backup_overlay.go    # Backup browser and restore overlay
│   │   ├── tagpicker_overlay.go # Tag filter overlay for report view
│   │   ├── idle.go              # Idle timer warning badge logic
│   │   ├── themes.go            # Theme definitions and ApplyTheme()
│   │   ├── styles.go            # Global style variables
│   │   └── keys.go              # Keybinding definitions
│   └── util/
│       └── format.go            # Duration formatting helpers
```

---

## License

[MIT](LICENSE)
