package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EstimateOptions are the estimate labels.
type EstimateOptions struct{ Points []int }

// Estimate selects story point chips; cue Select.
func Estimate(c *ui.Context, points *int, o EstimateOptions) ui.Element {
	if points == nil || *points < 0 {
		panic("project: Estimate needs nonnegative points")
	}
	if len(o.Points) == 0 {
		o.Points = []int{1, 2, 3, 5, 8, 13}
	}
	for _, value := range o.Points {
		if value <= 0 {
			panic("project: Estimate options must be positive")
		}
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, value := range o.Points {
			b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprint(value)).Padding(t.Space(1), t.Space(2)).Radius(999).
				Border(1, t.Border).Checked(*points == value))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*points == value)), ui.Snappy))
			b.Children(func() { ui.Text(c, fmt.Sprint(value)).FontSize(t.Rem(0.75)) })
			if b.Clicked() && *points != value {
				*points = value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Estimate", "selected", "points", value)
			}
		}
	})
	return root
}
