package project

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MilestoneOptions word the milestone; Percent takes a whole percent.
type MilestoneOptions struct{ Date, Percent, Done string }

// Milestone shows a diamond with date and progress.
func Milestone(c *ui.Context, title, date string, value float64, o MilestoneOptions) ui.Element {
	if title == "" || date == "" || value < 0 || value > 1 || !finite(value) {
		panic(fmt.Sprintf("project: Milestone %q has invalid data", title))
	}
	def(&o.Date, "Due %s")
	def(&o.Percent, "%d%%")
	def(&o.Done, "Reached")
	t := c.Theme()
	reached := value == 1
	ink := t.Accent
	if reached {
		ink = t.Success
	}
	root := card(c).Gap(t.Space(2.5)).Padding(t.Space(3))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Children(func() {
				gem := ui.Box(c).Size(t.Space(5.5), t.Space(5.5)).Rotate(45).Radius(t.Space(1)).
					Background(root.AnimateColor("gem", ink, ui.Snappy))
				kit.Arrive(gem, ui.Pose{Scale: 0.4, Rotate: -90, Fade: 1}, 0)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, title).FontSize(t.Rem(0.9)).Bold().SingleLine().MinWidth(0)
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "flag", 12)
					ui.Text(c, fmt.Sprintf(o.Date, date)).FontSize(t.Rem(0.75)).Shrink(0)
				})
			})
			if reached {
				c.Key("done")
				kit.Arrive(kit.Chip(c, o.Done, t.Success), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			} else {
				kit.Numeric(c, fmt.Sprintf(o.Percent, int(math.Floor(value*100)))).FontSize(t.Rem(1)).Bold().TextColor(ink)
			}
		})
		progress(c, value, ink)
	})
	return root
}
