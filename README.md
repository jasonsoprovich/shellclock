# shellclock

A terminal-based time tracker with a clean TUI, inspired by Toggl. Organize work into projects and tasks, run a live timer, and review time totals — all from the command line without leaving your terminal.

Built with the [Charm](https://charm.sh) stack: Bubble Tea, Lip Gloss, and Bubbles.

---

## Features

- **Projects & Tasks** — hierarchical tree: projects contain tasks, tasks accumulate sessions
- **Live Timer** — start, pause, resume, reset, and stop a timer on any task; timer persists across restarts
- **Session Editor** — view all recorded sessions for a task; add, edit, or delete individual sessions with second-level precision
- **Report View** — visual summary with progress bars showing time distribution across all projects and tasks
- **Theme System** — 31 built-in themes (Catppuccin family, Dracula, Nord, Tokyo Night, Gruvbox, Rosé Pine, and more); live preview in the picker; selection persists across launches
- **Keyboard-driven** — vim-style navigation throughout; context-sensitive help bar always visible
- **No external dependencies** — single JSON file, no database, no network

---

## Installation

### Prerequisites

- Go 1.21 or later

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

### Views

shellclock has four views, each accessible from the main tree.

---

### Tree View (main screen)

The collapsible project/task list. This is the starting point for all other views.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `→` / `l` | Expand project |
| `←` / `h` | Collapse project |
| `N` | New project |
| `n` | New task (under focused project) |
| `d` | Delete focused project or task |
| `enter` | Start timer on focused task |
| `e` | Open session editor for focused task |
| `R` | Open report view |
| `T` | Open theme picker |
| `q` | Quit |
| `?` | Toggle full key list |

---

### Timer View

Live timer for the active task. The timer keeps running if you navigate away.

| Key | Action |
|-----|--------|
| `p` | Pause / resume |
| `S` | Stop and save session |
| `r` | Reset (discard accumulated time) |
| `q` / `esc` | Return to tree (timer keeps running) |

The timer display shows:
- Current elapsed time (`HH:MM:SS`)
- Running / Paused badge
- Project › Task breadcrumb
- Wall-clock start time
- Previous total for the task (if any)

---

### Report View

Summary of all projects and tasks with time totals and proportional progress bars.

| Key | Action |
|-----|--------|
| `↑` / `k` | Scroll up |
| `↓` / `j` | Scroll down |
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

## Project structure

```
shellclock/
├── main.go                  # Entry point
├── internal/
│   ├── model/
│   │   └── model.go         # Data types and JSON persistence
│   ├── ui/
│   │   ├── app.go           # Root Bubble Tea model, view routing
│   │   ├── tree.go          # Project/task tree view
│   │   ├── timer.go         # Live timer view
│   │   ├── report.go        # Summary report view
│   │   ├── edit.go          # Session editor view
│   │   ├── themepicker.go   # Theme selection view
│   │   ├── themes.go        # Theme definitions and ApplyTheme()
│   │   ├── styles.go        # Global style variables
│   │   └── keys.go          # Keybinding definitions
│   └── util/
│       └── format.go        # Duration formatting helpers
```

---

## License

MIT
