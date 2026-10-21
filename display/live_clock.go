package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// LiveClock is the current time with seconds ticking.
func LiveClock(c *ui.Context) ui.Element {
	locale := c.Region().FormatLocale(c.Locale())
	text := i18n.FormatDate(locale, c.Now(), i18n.DateOptions{Time: i18n.MediumStyle})
	next := time.Second - time.Duration(c.Now().Nanosecond())
	if next <= 0 {
		next = time.Second
	}
	c.After(next)
	return kit.Numeric(c, text).FontFeatures("tnum")
}
