package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	wsTemperatureAsFahrenheit bool
	wsDistanceAsMiles         bool
	wsPrecipAsInches          bool
	wsWindAsMph               bool
)

var websocketCmd = &cobra.Command{
	Use:   "websocket",
	Short: "Live weather dashboard via WebSocket",
	Long: `Connect to the Tempest WebSocket API for a real-time weather dashboard.
Displays live observations (~60s), rapid wind (~3s), and weather events
(lightning, rain) in a full-screen terminal UI.

Press q to quit.`,
	Run: func(cmd *cobra.Command, args []string) {
		sid := cmd.Flag("station").Value.String()
		if sid == "" {
			fmt.Println("Station ID is required. Use -s <station_id>")
			return
		}

		apiToken := getAPIToken()

		// Fetch station info to get name, timezone, device ID
		station, err := fetchStationInfo(apiToken, sid)
		if err != nil {
			fmt.Printf("Error fetching station info: %v\n", err)
			return
		}

		if len(station.Stations) == 0 {
			fmt.Println("No station found for the given ID")
			return
		}

		s := station.Stations[0]
		deviceID := extractTempestDeviceID(station)
		if deviceID == 0 {
			fmt.Println("No Tempest device found for this station")
			return
		}

		model := DashboardModel{
			stationName: s.Name,
			timezone:    s.Timezone,
			deviceID:    deviceID,
			apiToken:    apiToken,
			unitPrefs: UnitPrefs{
				TempF:    wsTemperatureAsFahrenheit,
				WindMph:  wsWindAsMph,
				PrecipIn: wsPrecipAsInches,
				DistMi:   wsDistanceAsMiles,
			},
			wsDone: make(chan struct{}),
			msgCh:  make(chan tea.Msg, 32),
		}

		p := tea.NewProgram(model, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running dashboard: %v\n", err)
		}
	},
}

func fetchStationInfo(token, stationID string) (*Station, error) {
	baseURL := "https://swd.weatherflow.com/swd/rest/stations/" + stationID
	params := url.Values{}
	params.Add("token", token)

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.URL.RawQuery = params.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var station Station
	if err := json.Unmarshal(body, &station); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return &station, nil
}

func extractTempestDeviceID(station *Station) int {
	if len(station.Stations) == 0 {
		return 0
	}
	// Look for a Tempest device (device_type "ST")
	for _, device := range station.Stations[0].Devices {
		if device.DeviceType == "ST" {
			return device.DeviceID
		}
	}
	// Fallback to first device
	if len(station.Stations[0].Devices) > 0 {
		return station.Stations[0].Devices[0].DeviceID
	}
	return 0
}

func init() {
	rootCmd.AddCommand(websocketCmd)

	websocketCmd.Flags().BoolVarP(&wsTemperatureAsFahrenheit, "fahrenheit", "f", false, "Display temperature in Fahrenheit")
	websocketCmd.Flags().BoolVarP(&wsDistanceAsMiles, "miles", "", false, "Display distance in Miles")
	websocketCmd.Flags().BoolVarP(&wsPrecipAsInches, "inches", "", false, "Display precipitation in Inches")
	websocketCmd.Flags().BoolVarP(&wsWindAsMph, "mph", "", false, "Display wind speed in MPH")
}
