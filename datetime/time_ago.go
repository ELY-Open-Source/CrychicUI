package datetime

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// TimeAgoOptions tune a TimeAgo.
type TimeAgoOptions struct {
	// Now covers under a minute; Minutes, Hours, Days format counts.
	Now, Minutes, Hours, Days string
}

// TimeAgo says how long ago at was, updating on time.
func TimeAgo(c *ui.Context, at time.Time, o TimeAgoOptions) ui.Element {
	if at.IsZero() {
		panic("datetime: TimeAgo needs a time")
	}
	def(&o.Now, "just now")
	def(&o.Minutes, "%d min ago")
	def(&o.Hours, "%d h ago")
	def(&o.Days, "%d d ago")
	loc := locale(c)
	ago := max(c.Now().Sub(at), 0)
	var text string
	var step time.Duration
	switch {
	case ago < time.Minute:
		text, step = o.Now, time.Minute
	case ago < time.Hour:
		text, step = fmt.Sprintf(o.Minutes, int(ago/time.Minute)), time.Minute
	case ago < 24*time.Hour:
		text, step = fmt.Sprintf(o.Hours, int(ago/time.Hour)), time.Hour
	case ago < 7*24*time.Hour:
		text, step = fmt.Sprintf(o.Days, int(ago/(24*time.Hour))), 24*time.Hour
	default:
		text = i18n.FormatDate(loc, at, i18n.DateOptions{Skeleton: "yMMMd"})
	}
	if step > 0 {
		c.After(step - ago%step)
	}
	return ui.Text(c, text).ContentTransition(ui.ContentTransition{}).
		Tooltip(i18n.FormatDate(loc, at, i18n.DateOptions{Skeleton: "yMMMdjm"}))
}
