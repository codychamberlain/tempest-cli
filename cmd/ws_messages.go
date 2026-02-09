package cmd

import "time"

// --- Outbound WS messages ---

type WSListenStart struct {
	Type     string `json:"type"`
	DeviceID int    `json:"device_id"`
	ID       string `json:"id"`
}

type WSListenRapidStart struct {
	Type     string `json:"type"`
	DeviceID int    `json:"device_id"`
	ID       string `json:"id"`
}

// --- Inbound WS messages ---

// WSMessage is the generic envelope for type routing.
type WSMessage struct {
	Type string `json:"type"`
}

// WSObservation is an obs_st message from the Tempest WS.
type WSObservation struct {
	Type     string      `json:"type"`
	DeviceID int         `json:"device_id"`
	Obs      [][]float64 `json:"obs"`
}

// WSRapidWind is a rapid_wind message from the Tempest WS.
type WSRapidWind struct {
	Type     string    `json:"type"`
	DeviceID int       `json:"device_id"`
	Ob       []float64 `json:"ob"`
}

// WSEventPrecip is an evt_precip message from the Tempest WS.
type WSEventPrecip struct {
	Type     string    `json:"type"`
	DeviceID int       `json:"device_id"`
	Evt      []float64 `json:"evt"`
}

// WSEventStrike is an evt_strike message from the Tempest WS.
type WSEventStrike struct {
	Type     string    `json:"type"`
	DeviceID int       `json:"device_id"`
	Evt      []float64 `json:"evt"`
}

// WSAck is an acknowledgment message from the Tempest WS.
type WSAck struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// --- Parsed domain types ---

// ObsData holds named fields parsed from the obs_st 22-element array.
type ObsData struct {
	Timestamp      time.Time
	WindLull       float64
	WindAvg        float64
	WindGust       float64
	WindDirection  int
	Pressure       float64
	Temperature    float64
	Humidity       float64
	Illuminance    float64
	UV             float64
	SolarRadiation float64
	PrecipAccum    float64
	PrecipType     int
	LightningDist  float64
	LightningCount int
	Battery        float64
	DailyRain      float64
}

// RapidWindData holds parsed rapid_wind data.
type RapidWindData struct {
	Timestamp     time.Time
	WindSpeed     float64
	WindDirection int
}

// EventData represents a weather event (lightning, rain start).
type EventData struct {
	Timestamp time.Time
	Type      string
	Detail    string
}

// ParseObsArray converts the obs_st float64 array to named fields.
// obs_st indices per Tempest docs:
//
//	0:timestamp 1:wind_lull 2:wind_avg 3:wind_gust 4:wind_dir
//	5:wind_interval 6:pressure 7:temperature 8:humidity 9:illuminance
//	10:uv 11:solar_radiation 12:precip_accum 13:precip_type
//	14:lightning_distance 15:lightning_count 16:battery 17:report_interval
//	18:local_daily_rain 19-21:reserved
func ParseObsArray(obs []float64) ObsData {
	get := func(i int) float64 {
		if i < len(obs) {
			return obs[i]
		}
		return 0
	}
	return ObsData{
		Timestamp:      time.Unix(int64(get(0)), 0),
		WindLull:       get(1),
		WindAvg:        get(2),
		WindGust:       get(3),
		WindDirection:  int(get(4)),
		Pressure:       get(6),
		Temperature:    get(7),
		Humidity:       get(8),
		Illuminance:    get(9),
		UV:             get(10),
		SolarRadiation: get(11),
		PrecipAccum:    get(12),
		PrecipType:     int(get(13)),
		LightningDist:  get(14),
		LightningCount: int(get(15)),
		Battery:        get(16),
		DailyRain:      get(18),
	}
}

// ParseRapidWind converts the rapid_wind 3-element array to named fields.
// rapid_wind: [timestamp, wind_speed_m/s, wind_direction_deg]
func ParseRapidWind(ob []float64) RapidWindData {
	get := func(i int) float64 {
		if i < len(ob) {
			return ob[i]
		}
		return 0
	}
	return RapidWindData{
		Timestamp:     time.Unix(int64(get(0)), 0),
		WindSpeed:     get(1),
		WindDirection: int(get(2)),
	}
}
