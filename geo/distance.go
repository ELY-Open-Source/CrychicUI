package geo

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DistanceOptions dress the badge; zero Minutes hides the time.
type DistanceOptions struct {
	Icon, Time, Hours string
	Minutes           int
	Imperial          bool
	Units             Units
}

// Distance is a pill with distance and optional time.
func Distance(c *ui.Context, m float64, o DistanceOptions) ui.Element {
	def(&o.Icon, "navigation")
	def(&o.Time, "%d min")
	def(&o.Hours, "%d h %d min")
	o.Units.fill()
	meters("Distance", o.Icon, m)
	if o.Minutes < 0 {
		panic(fmt.Sprintf("geo: Distance %q needs minutes >= 0, got %d", o.Icon, o.Minutes))
	}
	t := c.Theme()
	root := ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2.5)).Radius(999).
		Background(t.Accent.Alpha(0.1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.9, Fade: 1}, 0)
	return root.Children(func() {
		glyph(c, o.Icon, 14).TextColor(t.Accent)
		kit.Numeric(c, span(m, o.Imperial, o.Units)).FontSize(t.Rem(0.8)).Bold().TextColor(t.Accent).Shrink(0).MaxWidth(t.Space(28)).SingleLine()
		if o.Minutes > 0 {
			c.Key("time")
			ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.Accent.Alpha(0.5))
			kit.Numeric(c, duration(o.Minutes, o.Time, o.Hours)).FontSize(t.Rem(0.78)).TextColor(t.Text.Mix(t.Accent, 0.4)).Shrink(0).MaxWidth(t.Space(28)).SingleLine()
		}
	})
}
