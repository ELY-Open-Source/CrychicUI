package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NextEventOptions word the card; In takes the countdown.
type NextEventOptions struct{ Label, In, Started, DateFormat string }

// NextEvent shows the next event with countdown.
func NextEvent(c *ui.Context, event Event, o NextEventOptions) ui.Element {
	validEvent("NextEvent", event)
	def(&o.Label, "Next event")
	def(&o.In, "in %s")
	def(&o.Started, "Happening now")
	def(&o.DateFormat, "Mon, Jan 2")
	t := c.Theme()
	ink := inks(t)[event.Tint]
	now := c.Now()
	// Whole seconds from Unix times: Sub saturates past 292 years.
	secs := event.Start.Unix() - now.Unix()
	if event.Start.Nanosecond() > now.Nanosecond() {
		secs++
	}
	if secs > 0 {
		frac := time.Duration((event.Start.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
		if frac == 0 {
			frac = time.Second
		}
		c.After(frac + time.Millisecond)
	}
	root := kit.Card(c).FillWidth().Row().AlignItems(ui.Stretch).Clip().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Width(t.Space(1.5)).Shrink(0).Background(ink)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink).Children(func() {
				glyph(c, "calendar-clock", 14)
				ui.Text(c, o.Label).FontSize(t.Rem(0.7)).Bold().Grow(1).MinWidth(0).SingleLine()
			})
			ui.Text(c, event.Title).FontSize(t.Rem(1.1)).Bold().SingleLine().MinWidth(0)
			meta(c, "clock", event.Start.Format(o.DateFormat)+" · "+eventText(event, event.Start.Location()))
			if event.Place != "" {
				meta(c, "map-pin", event.Place)
			}
			c.Key(secs > 0)
			if secs > 0 {
				kit.Numeric(c, fmt.Sprintf(o.In, countdown(secs))).FontSize(t.Rem(1.4)).Bold().TextColor(ink)
			} else {
				kit.Arrive(kit.Chip(c, o.Started, t.Success), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			}
		})
	})
	return root
}

// countdown is d hh:mm:ss, the days only when present.
func countdown(secs int64) string {
	clock := fmt.Sprintf("%02d:%02d:%02d", secs/3600%24, secs/60%60, secs%60)
	if days := secs / 86400; days > 0 {
		return fmt.Sprintf("%dd %s", days, clock)
	}
	return clock
}
