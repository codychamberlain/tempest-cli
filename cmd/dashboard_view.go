package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// --- Unit conversion (WS sends raw metric) ---

func convertTemp(c float64, prefs UnitPrefs) float64 {
	if prefs.TempF {
		return c*9.0/5.0 + 32.0
	}
	return c
}

func convertWind(ms float64, prefs UnitPrefs) float64 {
	if prefs.WindMph {
		return ms * 2.23694
	}
	return ms
}

func convertPrecip(mm float64, prefs UnitPrefs) float64 {
	if prefs.PrecipIn {
		return mm / 25.4
	}
	return mm
}

func windUnitLabel(prefs UnitPrefs) string {
	if prefs.WindMph {
		return "mph"
	}
	return "m/s"
}

func tempUnitLabel(prefs UnitPrefs) string {
	if prefs.TempF {
		return "F"
	}
	return "C"
}

func precipUnitLabel(prefs UnitPrefs) string {
	if prefs.PrecipIn {
		return "in"
	}
	return "mm"
}

// --- Icon inference (WS obs_st has no conditions string) ---

func inferWeatherIcon(obs *ObsData) string {
	if obs == nil {
		return "cloudy"
	}
	if obs.LightningCount > 0 {
		return "thunderstorm"
	}
	if obs.PrecipAccum > 0 && obs.Temperature <= 0 {
		return "snow"
	}
	if obs.PrecipAccum > 0 {
		return "rainy"
	}
	if obs.SolarRadiation < 10 && obs.Illuminance < 100 {
		return "clear-night"
	}
	if obs.Illuminance > 50000 {
		return "clear-day"
	}
	if obs.SolarRadiation > 200 {
		return "partly-cloudy-day"
	}
	return "cloudy"
}

// --- Direction helpers ---

func degreesToCardinal(deg int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	idx := int(math.Round(float64(deg)/45.0)) % 8
	if idx < 0 {
		idx += 8
	}
	return dirs[idx]
}

// --- Wind compass ---

func renderWindCompass(direction int) string {
	cardinal := degreesToCardinal(direction)

	ne, nw, se, sw := " ", " ", " ", " "
	e, w, n, s := " ", " ", " ", " "

	switch cardinal {
	case "N":
		n = "o"
	case "NE":
		ne = "o"
	case "E":
		e = "o"
	case "SE":
		se = "o"
	case "S":
		s = "o"
	case "SW":
		sw = "o"
	case "W":
		w = "o"
	case "NW":
		nw = "o"
	}

	_ = n // N is always labeled; marker shown via position

	return strings.Join([]string{
		fmt.Sprintf("       N       "),
		fmt.Sprintf("    %s  |  %s    ", nw, ne),
		fmt.Sprintf(" W %s--+---%s E ", w, e),
		fmt.Sprintf("    %s  |  %s    ", sw, se),
		fmt.Sprintf("       %s       ", s),
	}, "\n")
}

// --- Wind sparkline ---

func renderWindSparkline(history []RapidWindData) string {
	if len(history) == 0 {
		return ""
	}
	blocks := []rune("▁▂▃▄▅▆▇█")
	maxSpeed := 0.0
	minSpeed := math.MaxFloat64
	for _, w := range history {
		if w.WindSpeed > maxSpeed {
			maxSpeed = w.WindSpeed
		}
		if w.WindSpeed < minSpeed {
			minSpeed = w.WindSpeed
		}
	}
	spread := maxSpeed - minSpeed
	if spread < 0.1 {
		spread = 1
	}

	var sb strings.Builder
	for _, w := range history {
		idx := int((w.WindSpeed - minSpeed) / spread * float64(len(blocks)-1))
		if idx >= len(blocks) {
			idx = len(blocks) - 1
		}
		if idx < 0 {
			idx = 0
		}
		sb.WriteRune(blocks[idx])
	}
	return sb.String()
}

// --- Panel rendering ---

func renderDashboard(m DashboardModel) string {
	width := m.width
	if width <= 0 || width > 80 {
		width = 80
	}

	header := renderDashHeader(m, width)
	conditions := renderConditionsPanel(m, width)
	wind := renderWindPanel(m, width)
	events := renderEventsPanel(m, width)
	status := renderStatusBar(m, width)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		conditions,
		wind,
		events,
		status,
	)
}

func renderDashHeader(m DashboardModel, width int) string {
	theme := getWeatherTheme(inferWeatherIcon(m.currentObs))

	indicator := "[CONNECTING...]"
	indicatorColor := "#FFAA00"
	if m.connected {
		indicator = "[LIVE]"
		indicatorColor = "#00FF00"
	}
	if m.errMsg != "" && !m.reconnecting {
		indicator = "[DISCONNECTED]"
		indicatorColor = "#FF0000"
	}

	indicatorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(indicatorColor)).Bold(true)
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	innerWidth := width - 4 // account for border + padding
	title := fmt.Sprintf("%s (%s)", m.stationName, m.timezone)
	indStr := indicatorStyle.Render(indicator)
	titleStr := titleStyle.Render(title)

	// Pad to right-align indicator
	titleLen := lipgloss.Width(titleStr)
	indLen := lipgloss.Width(indStr)
	padding := innerWidth - titleLen - indLen
	if padding < 1 {
		padding = 1
	}

	content := titleStr + strings.Repeat(" ", padding) + indStr

	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(0, 1)

	return headerStyle.Render(content)
}

func renderConditionsPanel(m DashboardModel, width int) string {
	theme := getWeatherTheme(inferWeatherIcon(m.currentObs))

	panelStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(1, 2)

	if m.currentObs == nil {
		waitMsg := "Waiting for observation data..."
		if m.errMsg != "" {
			waitMsg = fmt.Sprintf("Error: %s", m.errMsg)
		}
		return panelStyle.Render(waitMsg)
	}

	obs := m.currentObs
	prefs := m.unitPrefs
	iconKey := inferWeatherIcon(obs)
	icon := getWeatherIcon(iconKey)

	iconStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary))

	iconBlock := iconStyle.Render(strings.Join(icon.Full, "\n"))

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Label))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	tempStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	temp := convertTemp(obs.Temperature, prefs)
	windSpeed := convertWind(obs.WindAvg, prefs)
	rain := convertPrecip(obs.DailyRain, prefs)
	dewPoint := convertTemp(obs.Temperature-(100-obs.Humidity)/5, prefs) // approximation

	condName := inferConditionName(obs)

	lines := []string{
		fmt.Sprintf("%s", tempStyle.Render(fmt.Sprintf("%.1f\u00b0%s", temp, tempUnitLabel(prefs)))),
		valueStyle.Render(condName),
		"",
		fmt.Sprintf("%s %s       %s %s",
			labelStyle.Render("Humidity:"),
			valueStyle.Render(fmt.Sprintf("%.0f%%", obs.Humidity)),
			labelStyle.Render("Wind:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s %s", windSpeed, windUnitLabel(prefs), degreesToCardinal(obs.WindDirection))),
		),
		fmt.Sprintf("%s %s      %s %s",
			labelStyle.Render("Pressure:"),
			valueStyle.Render(fmt.Sprintf("%.0f mb", obs.Pressure)),
			labelStyle.Render("UV:"),
			valueStyle.Render(fmt.Sprintf("%.0f", obs.UV)),
		),
		fmt.Sprintf("%s %s    %s %s",
			labelStyle.Render("Dew Point:"),
			valueStyle.Render(fmt.Sprintf("%.1f\u00b0%s", dewPoint, tempUnitLabel(prefs))),
			labelStyle.Render("Rain:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s", rain, precipUnitLabel(prefs))),
		),
	}

	stats := strings.Join(lines, "\n")
	content := lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "   ", stats)

	return panelStyle.Render(content)
}

func inferConditionName(obs *ObsData) string {
	if obs == nil {
		return "Unknown"
	}
	if obs.LightningCount > 0 {
		return "Thunderstorm"
	}
	if obs.PrecipAccum > 0 && obs.Temperature <= 0 {
		return "Snow"
	}
	if obs.PrecipAccum > 0 {
		return "Rain"
	}
	if obs.SolarRadiation < 10 && obs.Illuminance < 100 {
		return "Clear Night"
	}
	if obs.Illuminance > 50000 {
		return "Clear"
	}
	if obs.SolarRadiation > 200 {
		return "Partly Cloudy"
	}
	return "Cloudy"
}

func renderWindPanel(m DashboardModel, width int) string {
	theme := getWeatherTheme(inferWeatherIcon(m.currentObs))

	panelStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(0, 2)

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Label))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	prefs := m.unitPrefs

	windDir := 0
	windSpeed := 0.0
	windGust := 0.0
	windLull := 0.0
	if m.currentObs != nil {
		windDir = m.currentObs.WindDirection
		windSpeed = convertWind(m.currentObs.WindAvg, prefs)
		windGust = convertWind(m.currentObs.WindGust, prefs)
		windLull = convertWind(m.currentObs.WindLull, prefs)
	}
	// Override with rapid wind if available
	if m.rapidWind != nil {
		windDir = m.rapidWind.WindDirection
		windSpeed = convertWind(m.rapidWind.WindSpeed, prefs)
	}

	wUnit := windUnitLabel(prefs)

	compass := renderWindCompass(windDir)

	statsLines := []string{
		fmt.Sprintf("%s %s   %s %s",
			labelStyle.Render("Speed:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s", windSpeed, wUnit)),
			labelStyle.Render("Gust:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s", windGust, wUnit)),
		),
		fmt.Sprintf("%s  %s   %s  %s",
			labelStyle.Render("Lull:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s", windLull, wUnit)),
			labelStyle.Render("Dir:"),
			valueStyle.Render(fmt.Sprintf("%s %d\u00b0", degreesToCardinal(windDir), windDir)),
		),
		"",
		valueStyle.Render(renderWindSparkline(m.windHistory)),
	}

	compassBlock := lipgloss.NewStyle().Render(compass)
	statsBlock := strings.Join(statsLines, "\n")

	windContent := lipgloss.JoinHorizontal(lipgloss.Top, compassBlock, "    ", statsBlock)
	header := headerStyle.Render("  WIND")

	content := lipgloss.JoinVertical(lipgloss.Left, header, windContent)
	return panelStyle.Render(content)
}

func renderEventsPanel(m DashboardModel, width int) string {
	theme := getWeatherTheme(inferWeatherIcon(m.currentObs))

	panelStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(0, 2)

	headerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Label))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	header := headerStyle.Render("  EVENTS")

	if len(m.events) == 0 {
		content := lipgloss.JoinVertical(lipgloss.Left, header, labelStyle.Render("  No recent events"))
		return panelStyle.Render(content)
	}

	var lines []string
	for _, evt := range m.events {
		ts := evt.Timestamp.Format("15:04")
		line := fmt.Sprintf("  %s  %s", labelStyle.Render(ts), valueStyle.Render(evt.Detail))
		lines = append(lines, line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, append([]string{header}, lines...)...)
	return panelStyle.Render(content)
}

func renderStatusBar(m DashboardModel, width int) string {
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))

	var leftText string
	if m.lastUpdate.IsZero() {
		leftText = "  Waiting for data..."
	} else {
		ago := time.Since(m.lastUpdate).Truncate(time.Second)
		leftText = fmt.Sprintf("  Last updated: %s ago", ago)
	}

	if m.reconnecting {
		leftText += fmt.Sprintf(" (reconnecting %d/5...)", m.reconnectAttempts)
	}
	if m.errMsg != "" && !m.reconnecting {
		leftText += fmt.Sprintf(" | Error: %s", m.errMsg)
	}

	rightText := "Press q to quit  "

	innerWidth := width
	leftLen := lipgloss.Width(leftText)
	rightLen := lipgloss.Width(rightText)
	padding := innerWidth - leftLen - rightLen
	if padding < 1 {
		padding = 1
	}

	return labelStyle.Render(leftText + strings.Repeat(" ", padding) + rightText)
}
