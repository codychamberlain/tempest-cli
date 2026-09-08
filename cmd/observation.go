/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

type Observation struct {
	Status struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"status"`
	Elevation float64 `json:"elevation"`
	IsPublic  bool    `json:"is_public"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Obs       []struct {
		AirDensity                       float64 `json:"air_density"`
		AirTemperature                   float64 `json:"air_temperature"`
		BarometricPressure               float64 `json:"barometric_pressure"`
		Brightness                       int     `json:"brightness"`
		DeltaT                           float64 `json:"delta_t"`
		DewPoint                         float64 `json:"dew_point"`
		FeelsLike                        float64 `json:"feels_like"`
		HeatIndex                        float64 `json:"heat_index"`
		LightningStrikeCount             int     `json:"lightning_strike_count"`
		LightningStrikeCountLast1Hr      int     `json:"lightning_strike_count_last_1hr"`
		LightningStrikeCountLast3Hr      int     `json:"lightning_strike_count_last_3hr"`
		LightningStrikeLastDistance      int     `json:"lightning_strike_last_distance"`
		LightningStrikeLastEpoch         int     `json:"lightning_strike_last_epoch"`
		Precip                           float64 `json:"precip"`
		PrecipAccumLast1Hr               float64 `json:"precip_accum_last_1hr"`
		PrecipAccumLocalDay              float64 `json:"precip_accum_local_day"`
		PrecipAccumLocalDayFinal         float64 `json:"precip_accum_local_day_final"`
		PrecipAccumLocalYesterday        float64 `json:"precip_accum_local_yesterday"`
		PrecipAccumLocalYesterdayFinal   float64 `json:"precip_accum_local_yesterday_final"`
		PrecipAnalysisTypeYesterday      int     `json:"precip_analysis_type_yesterday"`
		PrecipMinutesLocalDay            int     `json:"precip_minutes_local_day"`
		PrecipMinutesLocalYesterday      int     `json:"precip_minutes_local_yesterday"`
		PrecipMinutesLocalYesterdayFinal int     `json:"precip_minutes_local_yesterday_final"`
		PressureTrend                    string  `json:"pressure_trend"`
		RelativeHumidity                 int     `json:"relative_humidity"`
		SeaLevelPressure                 float64 `json:"sea_level_pressure"`
		SolarRadiation                   int     `json:"solar_radiation"`
		StationPressure                  float64 `json:"station_pressure"`
		Timestamp                        int     `json:"timestamp"`
		Uv                               float64 `json:"uv"`
		WetBulbGlobeTemperature          float64 `json:"wet_bulb_globe_temperature"`
		WetBulbTemperature               float64 `json:"wet_bulb_temperature"`
		WindAvg                          float64 `json:"wind_avg"`
		WindChill                        float64 `json:"wind_chill"`
		WindDirection                    int     `json:"wind_direction"`
		WindGust                         float64 `json:"wind_gust"`
		WindLull                         float64 `json:"wind_lull"`
	} `json:"obs"`
	OutdoorKeys  []string `json:"outdoor_keys"`
	PublicName   string   `json:"public_name"`
	StationID    int      `json:"station_id"`
	StationName  string   `json:"station_name"`
	StationUnits struct {
		UnitsDirection string `json:"units_direction"`
		UnitsDistance  string `json:"units_distance"`
		UnitsOther     string `json:"units_other"`
		UnitsPrecip    string `json:"units_precip"`
		UnitsPressure  string `json:"units_pressure"`
		UnitsTemp      string `json:"units_temp"`
		UnitsWind      string `json:"units_wind"`
	} `json:"station_units"`
	Timezone string `json:"timezone"`
}

var (
	obsTemperatureAsFahrenheit bool
	obsDistanceAsMiles         bool
	obsPrecipAsInches          bool
	obsWindAsMph               bool
)

// observationCmd represents the observation command
var observationCmd = &cobra.Command{
	Use:   "observation",
	Short: "Show the latest observation from a station",
	Long: `Retrieve the most recent observation from a Tempest station.

Values default to metric units (°C, m/s, mb, mm, km) as returned by the
API. Use --fahrenheit, --mph, --inches, and --miles to change temperature, wind,
precipitation, and distance units, or -o JSON for the raw API response.`,
	Run: func(cmd *cobra.Command, args []string) {
		if (cmd.Flag("station").Value.String()) == "" {
			fmt.Println("Station ID is required")
		} else {
			apiToken := getAPIToken()
			sid := cmd.Flag("station").Value.String()
			baseURL := "https://swd.weatherflow.com/swd/rest/observations/station/" + sid
			params := url.Values{}
			params.Add("token", apiToken)
			req, err := http.NewRequest("GET", baseURL, nil)
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			req.URL.RawQuery = params.Encode()
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Error reading response body:", err)
				return
			}
			if cmd.Flag("output").Value.String() == "JSON" {
				fmt.Println(string(body))
			} else {
				var o Observation
				err = json.Unmarshal(body, &o)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					return
				}
				fmt.Println("*-----------------------------------------------*")
				fmt.Printf("| Station Name: %s\t\t|\n", o.StationName)
				fmt.Printf("| Public Name: %s\t\t|\n", o.PublicName)
				fmt.Printf("| Latitude: %.4f, Longitude: %.4f\t|\n", o.Latitude, o.Longitude)
				if len(o.Obs) > 0 {
					latestObs := o.Obs[len(o.Obs)-1]
					prefs := UnitPrefs{
						TempF:    obsTemperatureAsFahrenheit,
						WindMph:  obsWindAsMph,
						PrecipIn: obsPrecipAsInches,
						DistMi:   obsDistanceAsMiles,
					}
					tu := tempUnitLabel(prefs)
					wu := windUnitLabel(prefs)
					pu := precipUnitLabel(prefs)
					du := distUnitLabel(prefs)
					fmt.Println("|-----------------------------------------------|")
					fmt.Println("| Observation\t\t\t| Value\t\t|")
					fmt.Println("|-----------------------------------------------|")
					fmt.Printf("| Air Temperature\t\t| %.1f\u00B0%s\t|\n", convertTemp(latestObs.AirTemperature, prefs), tu)
					fmt.Printf("| Feels Like Temperature\t| %.1f\u00B0%s\t|\n", convertTemp(latestObs.FeelsLike, prefs), tu)
					fmt.Printf("| Dewpoint\t\t\t| %.1f\u00B0%s\t|\n", convertTemp(latestObs.DewPoint, prefs), tu)
					fmt.Printf("| Heat Index\t\t\t| %.1f\u00B0%s\t|\n", convertTemp(latestObs.HeatIndex, prefs), tu)
					fmt.Printf("| Wind Chill\t\t\t| %.1f\u00B0%s\t|\n", convertTemp(latestObs.WindChill, prefs), tu)
					fmt.Printf("| Relative Humidity\t\t| %d%%\t\t|\n", latestObs.RelativeHumidity)
					fmt.Printf("| Wind Direction\t\t| %d\u00B0\t\t|\n", latestObs.WindDirection)
					fmt.Printf("| Wind Speed\t\t\t| %.1f %s\t|\n", convertWind(latestObs.WindAvg, prefs), wu)
					fmt.Printf("| Air Density\t\t\t| %.2f\t\t|\n", latestObs.AirDensity)
					fmt.Printf("| Barometric Pressure\t\t| %.1f\t\t|\n", latestObs.BarometricPressure)
					fmt.Printf("| Pressure Trend\t\t| %s\t|\n", latestObs.PressureTrend)
					fmt.Printf("| Lightning Strike Count\t| %d\t\t|\n", latestObs.LightningStrikeCount)
					fmt.Printf("| Lightning Strike Distance\t| %.0f %s\t\t|\n", convertDist(float64(latestObs.LightningStrikeLastDistance), prefs), du)
					fmt.Printf("| Precip\t\t\t| %.2f %s\t|\n", convertPrecip(latestObs.Precip, prefs), pu)
					fmt.Printf("| Precip Last Hour\t\t| %.2f %s\t|\n", convertPrecip(latestObs.PrecipAccumLast1Hr, prefs), pu)
					fmt.Printf("| Precip Day\t\t\t| %.2f %s\t|\n", convertPrecip(latestObs.PrecipAccumLocalDayFinal, prefs), pu)
					fmt.Printf("| Precip Yesterday\t\t| %.2f %s\t|\n", convertPrecip(latestObs.PrecipAccumLocalYesterdayFinal, prefs), pu)
					fmt.Printf("| UV Index\t\t\t| %.1f\t\t|\n", latestObs.Uv)
					fmt.Printf("| Solar Radiation\t\t| %d W/m^2\t|\n", latestObs.SolarRadiation)
					fmt.Println("*-----------------------------------------------*")

				} else {
					fmt.Println("No observations available.")
				}
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(observationCmd)

	observationCmd.Flags().BoolVarP(&obsTemperatureAsFahrenheit, "fahrenheit", "f", false, "Display temperature in Fahrenheit (default: Celsius)")
	observationCmd.Flags().BoolVarP(&obsDistanceAsMiles, "miles", "", false, "Display distance in Miles (default: Kilometers)")
	observationCmd.Flags().BoolVarP(&obsPrecipAsInches, "inches", "", false, "Display precipitation in Inches (default: Millimeters)")
	observationCmd.Flags().BoolVarP(&obsWindAsMph, "mph", "", false, "Display wind speed in MPH (default: m/s)")
}
