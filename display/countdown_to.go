package display

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CountdownTo is days, hours, minutes and seconds to a time.
func CountdownTo(c *ui.Context, target time.Time) ui.Element {
	d := target.Sub(c.Now())
	if d < 0 {
		d = 0
	}
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second
	text := fmt.Sprintf("%dd %02dh %02dm %02ds", days, hours, minutes, seconds)
	next := d % time.Second
	if next == 0 {
		next = time.Second
	} else {
		next += time.Millisecond
	}
	c.After(next)
	return kit.Numeric(c, text).FontFeatures("tnum")
}
