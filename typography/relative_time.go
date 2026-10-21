package typography

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// RelativeTimeOptions override the English relative-time copy.
type RelativeTimeOptions struct {
	Now, Minute, Hour, Day string
	AgoSuffix, InPrefix    string
}

// RelativeTime is relative text updating with the clock.
func RelativeTime(c *ui.Context, when time.Time, o RelativeTimeOptions) ui.Element {
	if o.Now == "" {
		o.Now = "now"
	}
	if o.Minute == "" {
		o.Minute = "min"
	}
	if o.Hour == "" {
		o.Hour = "hr"
	}
	if o.Day == "" {
		o.Day = "d"
	}
	if o.AgoSuffix == "" {
		o.AgoSuffix = " ago"
	}
	if o.InPrefix == "" {
		o.InPrefix = "in "
	}
	now := c.Now()
	future := when.After(now)
	d := now.Sub(when)
	if d < 0 {
		d = -d
	}
	text := o.Now
	var unit time.Duration
	if d >= 24*time.Hour {
		unit = 24 * time.Hour
		text = fmt.Sprintf("%d %s", d/unit, o.Day)
	} else if d >= time.Hour {
		unit = time.Hour
		text = fmt.Sprintf("%d %s", int(d.Hours()), o.Hour)
	} else if d >= time.Minute {
		unit = time.Minute
		text = fmt.Sprintf("%d %s", int(d.Minutes()), o.Minute)
	}
	var next time.Duration
	switch {
	case d < time.Minute && future:
		next = d + time.Minute
	case d < time.Minute:
		next = time.Minute - d
	case future:
		next = d%unit + time.Millisecond
	default:
		next = unit - d%unit
	}
	if d >= time.Minute {
		if future {
			text = o.InPrefix + text
		} else {
			text += o.AgoSuffix
		}
	}
	c.After(next)
	return ui.Text(c, text).ContentTransition(ui.ContentTransition{})
}
