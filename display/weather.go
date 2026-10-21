package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WeatherCondition is the condition a Weather shows.
type WeatherCondition int

// Weather conditions.
const (
	WeatherSunny WeatherCondition = iota
	WeatherCloudy
	WeatherRainy
	WeatherSnowy
)

// Weather is a condition icon and temperature.
func Weather(c *ui.Context, condition WeatherCondition, temperature float64, unit string) ui.Element {
	if condition < WeatherSunny || condition > WeatherSnowy || !finite(temperature) || unit == "" {
		panic("display: Weather needs a condition, finite temperature and unit")
	}
	t := c.Theme()
	name := [...]string{"sun", "cloud", "cloud-rain", "snowflake"}[condition]
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		ui.Icon(c, icon.Must(name)).Size(t.Space(5), t.Space(5)).TextColor(t.Accent).ContentTransition(ui.ContentTransition{})
		kit.Numeric(c, formatNumber(c, temperature)+unit).Bold()
	})
}
