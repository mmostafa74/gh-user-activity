# GitHub User Activity CLI

<p>
  <img src="https://img.shields.io/badge/Go-1.26-00ADD8?logo=go" alt="Go version">
  <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey" alt="Platform">
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="License">
</p>

View recent GitHub user activity from the terminal. Fetch, format, and browse events for any GitHub user — with both a classic CLI mode and an interactive TUI built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Based on the [**GitHub User Activity**](https://roadmap.sh/projects/github-user-activity) project from [roadmap.sh](https://roadmap.sh).

---

## Table of Contents

- [Features](#features)
- [Demo](#demo)
- [Installation](#installation)
- [Dependencies](#dependencies)
- [Usage](#usage)
- [Examples](#examples)
- [Interactive TUI](#interactive-tui)
- [Data Storage](#data-storage)
- [Project Structure](#project-structure)
- [License](#license)

---

## Features

- View recent GitHub activity for any user
- Formatted event descriptions with timestamps
- Interactive TUI with scrollable event list
- Expand event details (commits, issue titles, PR titles)
- Plain-text event markers for terminal compatibility
- API call logging with file locking
- Supports 14+ GitHub event types
- Cross-platform

---

## Demo

![Demo](demo.gif)

---

## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.23 or later

### Steps

```bash
# Clone the repository
git clone https://github.com/yourusername/gh-user-activity.git
cd gh-user-activity

# Build the binary
go build -o gh-activity

# Verify it works
./gh-activity help
```

You can also run directly without building:

```bash
go run . view torvalds
```

---

## Dependencies

| Package | Purpose |
|---|---|
| [bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework |
| [bubbles](https://github.com/charmbracelet/bubbles) | Pre-built UI components (text input, spinner, viewport) |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal styling |
| [gofrs/flock](https://github.com/gofrs/flock) | Cross-platform file locking for log writes |

---

## Usage

| Command | Example | Description |
|---|---|---|
| `view` | `gh-activity view torvalds` | Show recent GitHub activity for a user |
| `tui` | `gh-activity tui` | Launch interactive TUI |
| `help` | `gh-activity help` | Show usage information |

---

## Examples

### View user activity

```bash
./gh-activity view torvalds
```

Sample output:

```
Activity for torvalds (last 30 events)
──────────────────────────────────────────────────
2026-07-06 14:32  • Pushed to master in torvalds/linux
2026-07-06 14:28  • Starred torvalds/linux
2026-07-06 14:25  • Opened issue #42 in torvalds/linux
2026-07-06 14:20  • Forked torvalds/linux
2026-07-06 14:15  • Created branch feature-x in torvalds/linux
```

### Error handling

```bash
./gh-activity view nonexistent-user    # Error: GitHub API returned 404 Not Found
./gh-activity                          # Shows usage
./gh-activity unknown                  # Shows usage
```

---

## Interactive TUI

The TUI mode provides a richer browsing experience:

```bash
./gh-activity tui
```

### Controls

| Key | Action |
|---|---|
| `enter` | Submit username / toggle event detail |
| `↑ / ↓` | Scroll through events |
| `esc` / `ctrl+r` | Back to search |
| `q` / `ctrl+c` | Quit |

### Screens

**Input screen** — enter a GitHub username:

```
         GitHub User Activity

    Enter a GitHub username:

    ┌────────────────────────────┐
    │ torvalds                   │
    └────────────────────────────┘

      [enter] submit  •  [q] quit
```

**Results screen** — browse events with scroll and expand:

```
   GitHub Activity: torvalds
──────────────────────────────────────────────────
[Push] Pushed to master in torvalds/linux
  2026-07-06 14:32
  commit a1b2c3d: Merge branch 'master'
  commit e4f5g6h: Fix memory leak in scheduler

[Star] Starred torvalds/linux
  2026-07-06 14:28

[Issue] Opened issue #42 in torvalds/linux
  2026-07-06 14:25
  #42: Kernel panic on ARM64 NUMA systems

 [↑/↓] scroll  [enter] toggle detail  [esc] search  [q] quit
```

---

## Data Storage

API calls are logged to `logs/apilogs.json` for debugging. The file is created and locked safely with `flock` when multiple requests write concurrently.

```json
[
  {
    "username": "torvalds",
    "status_code": 200,
    "response_body": "[{\"id\":\"...\",\"type\":\"PushEvent\",...}]"
  }
]
```

---

## Project Structure

| File | Responsibility |
|---|---|
| `main.go` | Entry point — CLI dispatch (`view`, `tui`, `help`) |
| `client/client.go` | GitHub API client — fetches user events |
| `client/events.go` | Event type definitions (Event, Actor, Repo, Payload) |
| `client/storage.go` | API call logging with `flock`-based file locking |
| `formatter/format.go` | Event formatting — human-readable descriptions |
| `tui/tui.go` | Interactive TUI with Bubble Tea |
| `utils/helpers.go` | Shared helpers (usage, event printing) |
| `go.mod` | Go module definition |

---

## License

MIT
