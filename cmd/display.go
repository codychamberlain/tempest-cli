package cmd

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// getTerminalWidth returns the terminal width, falling back to 100.
func getTerminalWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w < 40 {
		return 100
	}
	return w
}

// RenderForecast is the main entry point for the styled weather display.
func RenderForecast(f Forecast) {
	width := getTerminalWidth()
	if width > 80 {
		width = 80
	}

	theme := getWeatherTheme(f.CurrentConditions.Icon)

	header := renderHeader(f, width, theme)
	current := renderCurrentConditions(f, width, theme)
	daily := renderDailyForecast(f, width, theme)

	fmt.Println(header)
	fmt.Println(current)
	fmt.Println(daily)
}

// renderHeader creates the location/timezone banner.
func renderHeader(f Forecast, width int, theme WeatherTheme) string {
	headerStyle := lipgloss.NewStyle().
		Width(width - 2).
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(0, 1)

	title := fmt.Sprintf("%s (%s)", f.LocationName, f.Timezone)
	return headerStyle.Render(title)
}

// renderCurrentConditions creates the main weather panel with icon + stats.
func renderCurrentConditions(f Forecast, width int, theme WeatherTheme) string {
	icon := getWeatherIcon(f.CurrentConditions.Icon)

	iconStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Primary))

	iconBlock := iconStyle.Render(strings.Join(icon.Full, "\n"))

	stats := renderCurrentStats(f, theme)

	// Lay out icon on the left and stats on the right
	content := lipgloss.JoinHorizontal(lipgloss.Top, iconBlock, "   ", stats)

	panelStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(1, 2)

	return panelStyle.Render(content)
}

// renderCurrentStats builds the text block of weather statistics.
func renderCurrentStats(f Forecast, theme WeatherTheme) string {
	cc := f.CurrentConditions
	u := f.Units

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Label))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	tempStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	tempUnit := tempUnitSymbol(u.UnitsTemp)
	windUnit := u.UnitsWind

	lines := []string{
		fmt.Sprintf("%s  %s",
			tempStyle.Render(formatTemp(cc.AirTemperature, tempUnit)),
			labelStyle.Render(fmt.Sprintf("Feels like %s", formatTemp(cc.FeelsLike, tempUnit))),
		),
		valueStyle.Render(cc.Conditions),
		"",
		fmt.Sprintf("%s %s       %s %s",
			labelStyle.Render("Humidity:"),
			valueStyle.Render(fmt.Sprintf("%d%%", cc.RelativeHumidity)),
			labelStyle.Render("Wind:"),
			valueStyle.Render(fmt.Sprintf("%.1f %s %s", cc.WindAvg, windUnit, cc.WindDirectionCardinal)),
		),
		fmt.Sprintf("%s %s      %s %s",
			labelStyle.Render("Pressure:"),
			valueStyle.Render(formatPressure(cc.SeaLevelPressure, u.UnitsPressure)),
			labelStyle.Render("UV:"),
			valueStyle.Render(fmt.Sprintf("%d", cc.Uv)),
		),
		fmt.Sprintf("%s %s    %s %s",
			labelStyle.Render("Dew Point:"),
			valueStyle.Render(formatTemp(cc.DewPoint, tempUnit)),
			labelStyle.Render("Precip:"),
			valueStyle.Render(fmt.Sprintf("%d%%", cc.PrecipProbability)),
		),
	}

	return strings.Join(lines, "\n")
}

// renderDailyForecast creates the multi-day forecast panel.
func renderDailyForecast(f Forecast, width int, theme WeatherTheme) string {
	days := f.Forecast.Daily
	maxCards := 5
	if width < 60 {
		maxCards = 3
	}
	if len(days) < maxCards {
		maxCards = len(days)
	}
	// Skip today (index 0), show upcoming days
	start := 1
	if len(days) <= 1 {
		start = 0
	}
	end := start + maxCards
	if end > len(days) {
		end = len(days)
	}

	var cards []string
	for i := start; i < end; i++ {
		cards = append(cards, renderDailyCard(days[i], f.Units, theme))
	}

	content := lipgloss.JoinHorizontal(lipgloss.Top, cards...)

	panelStyle := lipgloss.NewStyle().
		Width(width - 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.Secondary)).
		Padding(0, 1).
		Align(lipgloss.Center)

	return panelStyle.Render(content)
}

// renderDailyCard builds a single day's forecast card.
func renderDailyCard(day ForecastDaily, units ForecastUnits, theme WeatherTheme) string {
	icon := getWeatherIcon(day.Icon)
	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(getWeatherTheme(day.Icon).Primary))
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Label))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	tempUnit := tempUnitSymbol(units.UnitsTemp)
	dayStr := dayName(day.DayStartLocal)

	miniArt := iconStyle.Render(strings.Join(icon.Mini, "\n"))

	temps := valueStyle.Render(fmt.Sprintf("%s/%s",
		formatTempShort(day.AirTempHigh, tempUnit),
		formatTempShort(day.AirTempLow, tempUnit),
	))

	condStr := day.Conditions
	if len(condStr) > 10 {
		condStr = condStr[:10]
	}
	conditions := labelStyle.Render(condStr)

	var precipLine string
	if day.PrecipProbability > 0 {
		precipLine = labelStyle.Render(fmt.Sprintf("%d%%", day.PrecipProbability))
	}

	cardStyle := lipgloss.NewStyle().
		Width(13).
		Align(lipgloss.Center).
		Padding(0, 1)

	content := strings.Join([]string{
		valueStyle.Render(dayStr),
		miniArt,
		temps,
		conditions,
		precipLine,
	}, "\n")

	return cardStyle.Render(content)
}

// --- Helpers ---

func tempUnitSymbol(unit string) string {
	switch unit {
	case "f":
		return "F"
	default:
		return "C"
	}
}

func formatTemp(temp float64, unit string) string {
	return fmt.Sprintf("%.0f°%s", math.Round(temp), unit)
}

func formatTempShort(temp float64, unit string) string {
	return fmt.Sprintf("%.0f°", math.Round(temp))
}

func formatPressure(p float64, unit string) string {
	return fmt.Sprintf("%.0f %s", p, unit)
}

func dayName(epoch int) string {
	t := time.Unix(int64(epoch), 0)
	return t.Format("Mon")
}
