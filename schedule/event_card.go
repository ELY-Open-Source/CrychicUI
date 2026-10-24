package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EventCardOptions format the event's date.
type EventCardOptions struct{ DateFormat string }

// EventCard shows title, time, place, and a color bar.
func EventCard(c *ui.Context, event Event, o EventCardOptions) ui.Element {
	validEvent("EventCard", event)
	def(&o.DateFormat, "Mon, Jan 2")
	t := c.Theme()
	ink := inks(t)[event.Tint]
	root := kit.Card(c).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Stretch).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Width(t.Space(1)).Shrink(0).Radius(999).Background(ink)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Text(c, event.Title).Bold().SingleLine().MinWidth(0)
			meta(c, "clock", event.Start.Format(o.DateFormat)+" · "+eventText(event, event.Start.Location()))
			if event.Place != "" {
				meta(c, "map-pin", event.Place)
			}
		})
	})
	return root
}

// meta is a muted icon and text line.
func meta(c *ui.Context, name, text string) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).MinWidth(0).TextColor(t.TextMuted).Children(func() {
		glyph(c, name, 12)
		ui.Text(c, text).FontSize(t.Rem(0.75)).SingleLine().Grow(1).MinWidth(0)
	})
}
