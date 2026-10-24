package project

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EstimateOptions word the picker; Unit takes the points.
type EstimateOptions struct {
	Points      []int
	Label, Unit string
}

// Estimate selects story point chips; cue Select.
func Estimate(c *ui.Context, points *int, o EstimateOptions) ui.Element {
	if points == nil || *points < 0 {
		panic("project: Estimate needs nonnegative points")
	}
	if len(o.Points) == 0 {
		o.Points = []int{1, 2, 3, 5, 8, 13}
	}
	for i, value := range o.Points {
		if value <= 0 || slices.Contains(o.Points[:i], value) {
			panic(fmt.Sprintf("project: Estimate option %d is not positive or repeats", value))
		}
	}
	def(&o.Label, "Story points")
	def(&o.Unit, "%d pts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			glyph(c, "gauge", 14).TextColor(t.TextMuted)
			ui.Text(c, o.Label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0)
			if *points > 0 {
				c.Key("value")
				kit.Arrive(kit.Chip(c, fmt.Sprintf(o.Unit, *points), t.Accent), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			}
		})
		row := ui.Row(c).Wrap().Gap(t.Space(1.5))
		row.Children(func() {
			for i, value := range o.Points {
				on := *points == value
				text := fmt.Sprint(value)
				c.Key(value)
				b := kit.Press(ui.ButtonBase(c).Label(text).Size(t.Space(9), t.Space(9)).Center().Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				b.TextColor(b.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
				kit.Stagger(c, row, b, i).Children(func() { kit.Numeric(c, text).FontSize(t.Rem(0.85)).Bold() })
				if b.Clicked() && !on {
					*points = value
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Estimate", "selected", "points", value)
				}
			}
		})
	})
	return root
}
