package cmd

// WeatherTheme holds colors for rendering a weather condition.
type WeatherTheme struct {
	Primary   string // main icon color
	Secondary string // border / accent color
	Label     string // stat label color
}

// WeatherIcon holds ASCII art in two sizes.
type WeatherIcon struct {
	Full []string // 7 lines, ~15 chars wide (current conditions)
	Mini []string // 3 lines, ~7 chars wide (daily cards)
}

var weatherThemes = map[string]WeatherTheme{
	"clear-day":              {Primary: "#FFD700", Secondary: "#FFD700", Label: "#888888"},
	"clear-night":            {Primary: "#C0C0C0", Secondary: "#7B68EE", Label: "#888888"},
	"partly-cloudy-day":      {Primary: "#FFD700", Secondary: "#A0A0A0", Label: "#888888"},
	"partly-cloudy-night":    {Primary: "#C0C0C0", Secondary: "#A0A0A0", Label: "#888888"},
	"cloudy":                 {Primary: "#A0A0A0", Secondary: "#A0A0A0", Label: "#888888"},
	"rainy":                  {Primary: "#4A90D9", Secondary: "#4A90D9", Label: "#888888"},
	"possibly-rainy-day":     {Primary: "#4A90D9", Secondary: "#4A90D9", Label: "#888888"},
	"possibly-rainy-night":   {Primary: "#4A90D9", Secondary: "#4A90D9", Label: "#888888"},
	"snow":                   {Primary: "#FFFFFF", Secondary: "#87CEEB", Label: "#888888"},
	"possibly-snow-day":      {Primary: "#FFFFFF", Secondary: "#87CEEB", Label: "#888888"},
	"possibly-snow-night":    {Primary: "#FFFFFF", Secondary: "#87CEEB", Label: "#888888"},
	"sleet":                  {Primary: "#87CEEB", Secondary: "#87CEEB", Label: "#888888"},
	"possibly-sleet-day":     {Primary: "#87CEEB", Secondary: "#87CEEB", Label: "#888888"},
	"possibly-sleet-night":   {Primary: "#87CEEB", Secondary: "#87CEEB", Label: "#888888"},
	"thunderstorm":           {Primary: "#FFD700", Secondary: "#9370DB", Label: "#888888"},
	"possibly-thunderstorm-day":   {Primary: "#FFD700", Secondary: "#9370DB", Label: "#888888"},
	"possibly-thunderstorm-night": {Primary: "#FFD700", Secondary: "#9370DB", Label: "#888888"},
	"windy":                  {Primary: "#87CEEB", Secondary: "#87CEEB", Label: "#888888"},
	"foggy":                  {Primary: "#D3D3D3", Secondary: "#D3D3D3", Label: "#888888"},
}

var weatherIcons = map[string]WeatherIcon{
	"clear-day": {
		Full: []string{
			`    \   /    `,
			`     .-.     `,
			`  - (   ) -  `,
			`     '-'     `,
			`    /   \    `,
		},
		Mini: []string{
			` \|/ `,
			` (o) `,
			` /|\ `,
		},
	},
	"clear-night": {
		Full: []string{
			`      .--.   `,
			`     /    )  `,
			`    |        `,
			`     \    )  `,
			`      '--'   `,
		},
		Mini: []string{
			`  _  `,
			` ( ) `,
			`  ~  `,
		},
	},
	"partly-cloudy-day": {
		Full: []string{
			`   \  /      `,
			`  _ /''.     `,
			`    \   )    `,
			`  /''.-'     `,
			`             `,
		},
		Mini: []string{
			` ~|/ `,
			` /o) `,
			` /|  `,
		},
	},
	"partly-cloudy-night": {
		Full: []string{
			`    .--.     `,
			`   / .-.)    `,
			`  | (        `,
			`   \ '-.)    `,
			`    '--'     `,
		},
		Mini: []string{
			` .-. `,
			`(  ) `,
			` '-' `,
		},
	},
	"cloudy": {
		Full: []string{
			`             `,
			`    .--.     `,
			` .-(    ).   `,
			`(___.__)__)  `,
			`             `,
		},
		Mini: []string{
			` .-. `,
			`(   )`,
			` '-' `,
		},
	},
	"rainy": {
		Full: []string{
			`    .--.     `,
			` .-(    ).   `,
			`(___.__)__)  `,
			` ' ' ' ' '  `,
			`' ' ' ' '   `,
		},
		Mini: []string{
			` .-. `,
			`(   )`,
			` '''`,
		},
	},
	"snow": {
		Full: []string{
			`    .--.     `,
			` .-(    ).   `,
			`(___.__)__)  `,
			` *  *  *  *  `,
			`*  *  *  *   `,
		},
		Mini: []string{
			` .-. `,
			`(   )`,
			` * * `,
		},
	},
	"sleet": {
		Full: []string{
			`    .--.     `,
			` .-(    ).   `,
			`(___.__)__)  `,
			` ' * ' * '   `,
			`* ' * ' *    `,
		},
		Mini: []string{
			` .-. `,
			`(   )`,
			` '* `,
		},
	},
	"thunderstorm": {
		Full: []string{
			`    .--.     `,
			` .-(    ).   `,
			`(___.__)__)  `,
			`  /_  /_     `,
			` /  /  /     `,
		},
		Mini: []string{
			` .-. `,
			`(   )`,
			` /_/ `,
		},
	},
	"windy": {
		Full: []string{
			`             `,
			` ~~~         `,
			`  ~~~~~~     `,
			` ~~~         `,
			`             `,
		},
		Mini: []string{
			`     `,
			` ~~~ `,
			`     `,
		},
	},
	"foggy": {
		Full: []string{
			`             `,
			` _ - _ - _ - `,
			` - _ - _ - _ `,
			` _ - _ - _ - `,
			`             `,
		},
		Mini: []string{
			` - - `,
			`- - -`,
			` - - `,
		},
	},
}

func getWeatherTheme(icon string) WeatherTheme {
	if t, ok := weatherThemes[icon]; ok {
		return t
	}
	return weatherThemes["cloudy"]
}

func getWeatherIcon(icon string) WeatherIcon {
	// Direct match
	if ic, ok := weatherIcons[icon]; ok {
		return ic
	}
	// Map variant icons to their base art
	aliases := map[string]string{
		"possibly-rainy-day":           "rainy",
		"possibly-rainy-night":         "rainy",
		"possibly-snow-day":            "snow",
		"possibly-snow-night":          "snow",
		"possibly-sleet-day":           "sleet",
		"possibly-sleet-night":         "sleet",
		"possibly-thunderstorm-day":    "thunderstorm",
		"possibly-thunderstorm-night":  "thunderstorm",
	}
	if base, ok := aliases[icon]; ok {
		return weatherIcons[base]
	}
	return weatherIcons["cloudy"]
}
