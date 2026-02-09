# CLAUDE.md

## Project Overview

tempest-cli is a Go CLI application for accessing WeatherFlow Tempest station and forecast data. It uses the Tempest REST API (`swd.weatherflow.com`).

## Build & Run

```bash
go build .
./tempest-cli forecast -s <station_id>
```

A `.env` file with `API_TOKEN=<token>` is required in the working directory.

## Architecture

- **Framework**: Cobra for CLI structure, Lipgloss for terminal styling
- **Entry point**: `main.go` -> `cmd.Execute()`
- **Commands**: Each command is a file in `cmd/` with its own `init()` that registers with `rootCmd`
- **API token**: Loaded from `.env` via `godotenv` in `getAPIToken()` (defined in `forecast.go`)
- **Output modes**: All commands check `-o JSON` flag for raw JSON, otherwise render styled output

## Key Files

| File | Purpose |
|---|---|
| `cmd/forecast.go` | Forecast command, API call, all forecast data types (`Forecast`, `ForecastCurrentConditions`, `ForecastDaily`, `ForecastHourly`, `ForecastUnits`) |
| `cmd/display.go` | Lipgloss rendering: `RenderForecast()` entry point, header/current/daily panels, formatting helpers |
| `cmd/weather_icons.go` | ASCII art maps (full-size + mini), `WeatherTheme` color definitions, icon lookup with alias fallback |
| `cmd/observation.go` | Observation command and `Observation` struct |
| `cmd/station.go` | Station command and `Station` struct |
| `cmd/root.go` | Root command, global flags (`--station`, `--output`) |

## Conventions

- Named types are prefixed with the parent context (e.g. `ForecastDaily`, `ForecastUnits`)
- Anonymous structs are used for small nested objects that don't need to be referenced elsewhere (e.g. `Station`, `Status` inside `Forecast`)
- JSON output bypasses all formatting — just prints the raw API response body
- Weather icons use pure ASCII art (no emoji) with color applied via Lipgloss styles
- Terminal width is detected via `golang.org/x/term`, capped at 80 cols, fallback to 100

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/charmbracelet/lipgloss` — terminal styling (borders, colors, layout)
- `golang.org/x/term` — terminal width detection
- `github.com/joho/godotenv` — `.env` file loading
- `github.com/gorilla/websocket` — WebSocket support (used by websocket command)

## Testing

No test suite yet. Manual testing:

```bash
./tempest-cli forecast -s <station_id> --fahrenheit --mph
./tempest-cli forecast -s <station_id> -o JSON
./tempest-cli observation -s <station_id>
./tempest-cli station -s <station_id>
```
