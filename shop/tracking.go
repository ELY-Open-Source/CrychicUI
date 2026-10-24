package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TrackEvent is one carrier event; Done ones come first.
type TrackEvent struct {
	Label, Location string
	Time            time.Time
	Done            bool
}

// TrackingOptions lay out event times.
type TrackingOptions struct{ Time string }

// Tracking shows carrier events on a rail, latest done ringed.
func Tracking(c *ui.Context, events []TrackEvent, o TrackingOptions) ui.Element {
	if len(events) == 0 {
		panic("shop: Tracking needs events")
	}
	current := -1
	for i, e := range events {
		if e.Label == "" || e.Location == "" || e.Done && i != current+1 {
			panic(fmt.Sprintf("shop: Tracking event %q needs text, a location, and done events first", e.Label))
		}
		if e.Done {
			current = i
		}
	}
	def(&o.Time, "Jan 2, 15:04")
	t := c.Theme()
	rail := func(i int) ui.Color {
		if i <= 0 || i >= len(events) {
			return t.Border.Alpha(0)
		}
		return t.Border.Mix(t.Success, unit(events[i].Done))
	}
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup)
	root.Children(func() {
		for i, e := range events {
			c.Key(i)
			on := i == current
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Column(c).Width(t.Space(5)).Shrink(0).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Width(t.Space(0.5)).Height(t.Space(1.5)).Background(rail(i))
					dot := ui.Box(c).Size(t.Space(3), t.Space(3)).Shrink(0).Radius(999)
					switch {
					case on:
						dot.Background(t.Accent).Border(t.Space(1)*max(0, dot.AnimateTo("ring", 1, ui.Bouncy)), t.Accent.Alpha(0.25)).
							Size(t.Space(4), t.Space(4))
					case e.Done:
						dot.Background(t.Success)
					default:
						dot.Background(t.Background).Border(1, t.Border)
					}
					ui.Box(c).Width(t.Space(0.5)).Grow(1).Background(rail(i + 1))
				})
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(0.5), 0, t.Space(4), 0).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Start).Children(func() {
						label := ui.Text(c, e.Label).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
						if on {
							label.Bold()
						} else if !e.Done {
							label.TextColor(t.TextMuted)
						}
						kit.Muted(c, e.Time.Format(o.Time)).FontSize(t.Rem(0.75)).Shrink(0)
					})
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
						ui.Icon(c, icon.Must("map-pin")).Size(12, 12)
						ui.Text(c, e.Location).FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
			})
		}
	})
	return root
}
