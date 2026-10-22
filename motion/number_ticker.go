package motion

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// NumberTickerOptions tune a NumberTicker.
type NumberTickerOptions struct {
	// From is where the first count starts.
	From float64
	// Format writes the number; "%.0f" when empty.
	Format string
}

var count = []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}

// NumberTicker counts to value with easing, from where it shows.
func NumberTicker(c *ui.Context, value float64, o NumberTickerOptions) ui.Element {
	finite("NumberTicker", value)
	finite("NumberTicker", o.From)
	if o.Format == "" {
		o.Format = "%.0f"
	}
	root := ui.Box(c)
	from := ui.Local(root, "from", func() float64 { return o.From })
	to := ui.Local(root, "to", func() float64 { return value })
	shown := ui.Local(root, "shown", func() float64 { return o.From })
	if value != *to {
		*from, *to = *shown, value
	}
	return root.Children(func() {
		w := ui.Box(c.Key(*to))
		*shown = *to
		if *from != *to {
			if k := w.Keyframes("crychic.count", count); k < 1 {
				*shown = *from + (*to-*from)*float64(k)
			}
		}
		w.Children(func() {
			ui.Text(c, fmt.Sprintf(o.Format, *shown)).FontFeatures("tnum").Label(fmt.Sprintf(o.Format, *to))
		})
	})
}
