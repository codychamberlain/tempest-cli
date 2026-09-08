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

- **Framework**: Cobra for CLI structure, Lipgloss for terminal styling, Bubble Tea for the live dashboard TUI (`websocket` command)
- **Entry point**: `main.go` -> `cmd.Execute()`
- **Commands**: Each command is a file in `cmd/` with its own `init()` that registers with `rootCmd`
- **API token**: Loaded from `.env` via `godotenv` in `getAPIToken()` (defined in `forecast.go`)
- **Output modes**: `forecast`, `observation`, and `station` check the `-o JSON` flag for raw JSON, otherwise render styled output. `websocket` is an interactive TUI and ignores `-o`.
- **Live dashboard**: `websocket` fetches station info via REST to find the Tempest device (`device_type == "ST"`), then dials `wss://ws.weatherflow.com/swd/data` and sends `listen_start` + `listen_rapid_start`. A read-loop goroutine parses `obs_st`, `rapid_wind`, `evt_precip`, and `evt_strike` messages and forwards them to the Bubble Tea model over a channel. Reconnects up to 5 times with a 5s delay; 11 minute read deadline detects stale connections. Press `q` to quit.

## Key Files

| File | Purpose |
|---|---|
| `cmd/forecast.go` | Forecast command, API call, all forecast data types (`Forecast`, `ForecastCurrentConditions`, `ForecastDaily`, `ForecastHourly`, `ForecastUnits`) |
| `cmd/display.go` | Lipgloss rendering: `RenderForecast()` entry point, header/current/daily panels, formatting helpers |
| `cmd/weather_icons.go` | ASCII art maps (full-size + mini), `WeatherTheme` color definitions, icon lookup with alias fallback |
| `cmd/observation.go` | Observation command and `Observation` struct |
| `cmd/station.go` | Station command and `Station` struct |
| `cmd/root.go` | Root command, global flags (`--station`, `--output`) |
| `cmd/websocket.go` | Websocket command, unit flags, station/device lookup (`fetchStationInfo`, `extractTempestDeviceID`) |
| `cmd/dashboard.go` | `DashboardModel` (Bubble Tea model), tea message types, WS connect/read loop, reconnect logic |
| `cmd/dashboard_view.go` | Dashboard rendering: header, conditions, wind (compass + sparkline), events, status bar; unit conversion and icon inference from raw obs |
| `cmd/ws_messages.go` | WS message structs (`WSObservation`, `WSRapidWind`, `WSEventPrecip`, `WSEventStrike`), parsed types (`ObsData`, `RapidWindData`, `EventData`), array parsers |

## Conventions

- Named types are prefixed with the parent context (e.g. `ForecastDaily`, `ForecastUnits`)
- Anonymous structs are used for small nested objects that don't need to be referenced elsewhere (e.g. `Station`, `Status` inside `Forecast`)
- JSON output bypasses all formatting — just prints the raw API response body
- Weather icons use pure ASCII art (no emoji) with color applied via Lipgloss styles
- WebSocket data arrives as raw metric arrays with no conditions string; `dashboard_view.go` converts units per `UnitPrefs` and infers an icon/condition name from illuminance, solar radiation, precip, and lightning fields
- The station observation REST endpoint ignores `units_*` query params and always returns metric; `observation` converts temp, wind, and precip client-side with the `convert*` helpers from `dashboard_view.go`. The `station_units` field in its response reflects the user's app preferences, not the units of the returned values
- Websocket unit flags mirror the forecast flags (`--fahrenheit`, `--miles`, `--inches`, `--mph`) but use `ws`-prefixed vars
- Terminal width is detected via `golang.org/x/term`, capped at 80 cols, fallback to 100

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/charmbracelet/lipgloss` — terminal styling (borders, colors, layout)
- `golang.org/x/term` — terminal width detection
- `github.com/joho/godotenv` — `.env` file loading
- `github.com/gorilla/websocket` — WebSocket client (used by websocket command)
- `github.com/charmbracelet/bubbletea` — TUI framework for the live dashboard

## Testing

No test suite yet. Manual testing:

```bash
./tempest-cli forecast -s <station_id> --fahrenheit --mph
./tempest-cli forecast -s <station_id> -o JSON
./tempest-cli observation -s <station_id> --fahrenheit --mph --inches
./tempest-cli station -s <station_id>
./tempest-cli websocket -s <station_id> --fahrenheit --mph
```
