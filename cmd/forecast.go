/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type ForecastCurrentConditions struct {
	AirDensity                      float64 `json:"air_density"`
	AirTemperature                  float64 `json:"air_temperature"`
	Brightness                      int     `json:"brightness"`
	Conditions                      string  `json:"conditions"`
	DeltaT                          float64 `json:"delta_t"`
	DewPoint                        float64 `json:"dew_point"`
	FeelsLike                       float64 `json:"feels_like"`
	Icon                            string  `json:"icon"`
	IsPrecipLocalDayRainCheck       bool    `json:"is_precip_local_day_rain_check"`
	IsPrecipLocalYesterdayRainCheck bool    `json:"is_precip_local_yesterday_rain_check"`
	LightningStrikeCountLast1Hr     int     `json:"lightning_strike_count_last_1hr"`
	LightningStrikeCountLast3Hr     int     `json:"lightning_strike_count_last_3hr"`
	LightningStrikeLastDistance     int     `json:"lightning_strike_last_distance"`
	LightningStrikeLastDistanceMsg  string  `json:"lightning_strike_last_distance_msg"`
	LightningStrikeLastEpoch        int     `json:"lightning_strike_last_epoch"`
	PrecipAccumLocalDay             int     `json:"precip_accum_local_day"`
	PrecipAccumLocalYesterday       int     `json:"precip_accum_local_yesterday"`
	PrecipMinutesLocalDay           int     `json:"precip_minutes_local_day"`
	PrecipMinutesLocalYesterday     int     `json:"precip_minutes_local_yesterday"`
	PrecipProbability               int     `json:"precip_probability"`
	PressureTrend                   string  `json:"pressure_trend"`
	RelativeHumidity                int     `json:"relative_humidity"`
	SeaLevelPressure                float64 `json:"sea_level_pressure"`
	SolarRadiation                  int     `json:"solar_radiation"`
	StationPressure                 float64 `json:"station_pressure"`
	Time                            int     `json:"time"`
	Uv                              int     `json:"uv"`
	WetBulbGlobeTemperature         float64 `json:"wet_bulb_globe_temperature"`
	WetBulbTemperature              float64 `json:"wet_bulb_temperature"`
	WindAvg                         float64 `json:"wind_avg"`
	WindDirection                   int     `json:"wind_direction"`
	WindDirectionCardinal           string  `json:"wind_direction_cardinal"`
	WindGust                        float64 `json:"wind_gust"`
}

type ForecastDaily struct {
	AirTempHigh       float64 `json:"air_temp_high"`
	AirTempLow        float64 `json:"air_temp_low"`
	Conditions        string  `json:"conditions"`
	DayNum            int     `json:"day_num"`
	DayStartLocal     int     `json:"day_start_local"`
	Icon              string  `json:"icon"`
	MonthNum          int     `json:"month_num"`
	PrecipIcon        string  `json:"precip_icon"`
	PrecipProbability int     `json:"precip_probability"`
	PrecipType        string  `json:"precip_type"`
	Sunrise           int     `json:"sunrise"`
	Sunset            int     `json:"sunset"`
}

type ForecastHourly struct {
	AirTemperature        float64 `json:"air_temperature"`
	Conditions            string  `json:"conditions"`
	FeelsLike             float64 `json:"feels_like"`
	Icon                  string  `json:"icon"`
	LocalDay              int     `json:"local_day"`
	LocalHour             int     `json:"local_hour"`
	Precip                float64 `json:"precip"`
	PrecipIcon            string  `json:"precip_icon"`
	PrecipProbability     int     `json:"precip_probability"`
	PrecipType            string  `json:"precip_type"`
	RelativeHumidity      int     `json:"relative_humidity"`
	SeaLevelPressure      float64 `json:"sea_level_pressure"`
	StationPressure       float64 `json:"station_pressure"`
	Time                  int     `json:"time"`
	Uv                    float64 `json:"uv"`
	WindAvg               float64 `json:"wind_avg"`
	WindDirection         int     `json:"wind_direction"`
	WindDirectionCardinal string  `json:"wind_direction_cardinal"`
	WindGust              float64 `json:"wind_gust"`
}

type ForecastUnits struct {
	UnitsAirDensity     string `json:"units_air_density"`
	UnitsBrightness     string `json:"units_brightness"`
	UnitsDistance       string `json:"units_distance"`
	UnitsOther          string `json:"units_other"`
	UnitsPrecip         string `json:"units_precip"`
	UnitsPressure       string `json:"units_pressure"`
	UnitsSolarRadiation string `json:"units_solar_radiation"`
	UnitsTemp           string `json:"units_temp"`
	UnitsWind           string `json:"units_wind"`
}

type Forecast struct {
	CurrentConditions ForecastCurrentConditions `json:"current_conditions"`
	Forecast          struct {
		Daily  []ForecastDaily  `json:"daily"`
		Hourly []ForecastHourly `json:"hourly"`
	} `json:"forecast"`
	Latitude           float64 `json:"latitude"`
	LocationName       string  `json:"location_name"`
	Longitude          float64 `json:"longitude"`
	SourceIDConditions int     `json:"source_id_conditions"`
	Station            struct {
		Agl             float64 `json:"agl"`
		Elevation       float64 `json:"elevation"`
		IsStationOnline bool    `json:"is_station_online"`
		State           int     `json:"state"`
		StationID       int     `json:"station_id"`
	} `json:"station"`
	Status struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
	} `json:"status"`
	Timezone              string        `json:"timezone"`
	TimezoneOffsetMinutes int           `json:"timezone_offset_minutes"`
	Units                 ForecastUnits `json:"units"`
}

var (
	temperatureAsFahrenheit bool
	distanceAsMiles         bool
	precipAsInches          bool
	windAsMph               bool
)

// forecastCmd represents the forecast command
var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Given a station ID, return forecoast",
	Long: `Forecast data is available for Tempest stations. Given a specific API token and station id you will be returned with
	forestcast data. `,
	Run: func(cmd *cobra.Command, args []string) {
		apiToken := getAPIToken()
		sid := cmd.Flag("station").Value.String()
		baseURL := "https://swd.weatherflow.com/swd/rest/better_forecast"
		params := url.Values{}
		params.Add("station_id", sid)
		params.Add("token", apiToken)
		if temperatureAsFahrenheit {
			params.Add("units_temp", "f")
		}
		if distanceAsMiles {
			params.Add("units_distance", "mi")
		}
		if precipAsInches {
			params.Add("units_precip", "in")
		}
		if windAsMph {
			params.Add("units_wind", "mph")
		}
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
			var f Forecast
			err := json.Unmarshal(body, &f)
			if err != nil {
				log.Fatalf("Error unmarshaling JSON: %s", err)
			}
			RenderForecast(f)
		}

	},
}

func getAPIToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	return os.Getenv("API_TOKEN")
}

func init() {
	rootCmd.AddCommand(forecastCmd)

	forecastCmd.Flags().BoolVarP(&temperatureAsFahrenheit, "fahrenheit", "f", false, "Display temperature in Fahrenheit, Default is Celsius")
	forecastCmd.Flags().BoolVarP(&distanceAsMiles, "miles", "", false, "Display distance in Miles, Default is Kilometers")
	forecastCmd.Flags().BoolVarP(&precipAsInches, "inches", "", false, "Display precip in Inches Default is Millimeters")
	forecastCmd.Flags().BoolVarP(&windAsMph, "mph", "", false, "Display wind in Miles per hour default is Kilometers per hour")
}
