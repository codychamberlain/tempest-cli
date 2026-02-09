# Tempest-CLI

CLI access to your [TempestWX](https://tempestwx.com/) station and forecast data.

Built with Go, [Cobra](https://github.com/spf13/cobra), and [Lipgloss](https://github.com/charmbracelet/lipgloss).

## Setup

1. Clone the repo and build:

```bash
go build .
```

2. Create a `.env` file in the project root with your Tempest API token:

```
API_TOKEN=your_token_here
```

You can get an API token from [tempestwx.com](https://tempestwx.com/).

## Usage

```
tempest-cli [command] [flags]
```

### Global Flags

| Flag | Short | Description |
|---|---|---|
| `--station` | `-s` | Station ID to pull data from |
| `--output` | `-o` | Output format (`JSON` for raw JSON) |

### Commands

#### forecast

Retrieve the weather forecast for a station with a styled ASCII art display showing current conditions and a 5-day outlook.

```bash
tempest-cli forecast -s <station_id>
```

```
╭──────────────────────────────────────────────────────────────────────────────╮
│                     Brookmore Victoria (America/Chicago)                     │
╰──────────────────────────────────────────────────────────────────────────────╯
╭──────────────────────────────────────────────────────────────────────────────╮
│                                                                              │
│     \  /         33°F  Feels like 28°F                                       │
│      .-.         Partly Cloudy                                               │
│   - (  .''-.                                                                 │
│      '(    )     Humidity: 57%       Wind: 5.0 mph W                         │
│      '''-''      Pressure: 1021 mb      UV: 2                                │
│                  Dew Point: 19°F    Precip: 0%                               │
│                                                                              │
╰──────────────────────────────────────────────────────────────────────────────╯
╭──────────────────────────────────────────────────────────────────────────────╮
│       Mon            Tue            Wed            Thu            Fri        │
│    \|/ .--.       \|/ .--.          .--.           .--.           .--.       │
│    (o)(   )       (o)(   )        -(    )-       -(    )-       -(    )-     │
│    /|\  '--       /|\  '--         * * *          * * *          * * *       │
│     41°/28°        35°/26°        33°/19°        37°/27°        37°/26°      │
│  Partly Cloudy  Partly Cloudy  Snow Possible   Snow Likely   Snow Possible   │
│       10%            10%            30%            50%            30%        │
╰──────────────────────────────────────────────────────────────────────────────╯
```

Forecast flags:

| Flag | Short | Description |
|---|---|---|
| `--fahrenheit` | `-f` | Display temperature in Fahrenheit (default: Celsius) |
| `--miles` | | Display distance in miles (default: km) |
| `--inches` | | Display precipitation in inches (default: mm) |
| `--mph` | | Display wind in mph (default: km/h) |

#### observation

Retrieve the latest observation data from a station.

```bash
tempest-cli observation -s <station_id>
```

#### station

List all stations or get details for a specific station.

```bash
tempest-cli station                # list all stations
tempest-cli station -s <station_id>  # details for a specific station
```

#### websocket

WebSocket connection to station (in development).

```bash
tempest-cli websocket
```

### JSON Output

Any command supports raw JSON output with `-o JSON`:

```bash
tempest-cli forecast -s <station_id> -o JSON
```

## Project Structure

```
cmd/
  root.go             # root command and global flags
  forecast.go         # forecast command, API call, data types
  display.go          # lipgloss-styled terminal rendering
  weather_icons.go    # ASCII art icons and color themes
  observation.go      # observation command
  station.go          # station command
  websocket.go        # websocket command (stub)
main.go               # entry point
```

## License

MIT - see [LICENSE](LICENSE).
