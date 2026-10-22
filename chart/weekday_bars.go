package chart

import "github.com/egoist/mygo/ui"

// WeekdayBarsOptions tune a WeekdayBars.
type WeekdayBarsOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Days names Monday to Sunday; English when empty.
	Days [7]string
}

// WeekdayBars sets each week's bar side by side per weekday.
func WeekdayBars(c *ui.Context, label string, weeks []Series, at *int, o WeekdayBarsOptions) ui.Element {
	english := [7]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	for i := range o.Days {
		def(&o.Days[i], english[i])
	}
	return grouped(c, "WeekdayBars", label, o.Days[:], weeks, at, LineOptions{Height: o.Height, Format: o.Format})
}
