package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Beacon is a pulsing hotspot; a click opens its tip.
func Beacon(c *ui.Context, label, tip string) ui.Element {
	if label == "" || tip == "" {
		panic("feedback: Beacon needs a label and a tip")
	}
	t := c.Theme()
	d := t.Space(3)
	root := ui.ButtonBase(c).Label(label).Size(d*2, d*2).Shrink(0)
	open := ui.Local(root, "open", func() bool { return false })
	if kit.Clicked(root, sound.Open) {
		*open = !*open
		kit.Log("Beacon", "toggled", "open", *open)
	}
	root.Children(func() {
		halo := ui.Box(c).Absolute().Top(0).Left(0).Right(0).Bottom(0).Radius(999)
		k := halo.Loop("pulse", 1600*time.Millisecond, ui.EaseOut)
		halo.Background(t.Accent.Alpha(0.4*(1-k))).Scale(0.5+0.7*k, 0.5+0.7*k)
		ui.Box(c).Size(d, d).Radius(999).Background(t.Accent).Border(t.Space(0.5), t.Background)
	})
	ui.PopoverBase(c, root, open, func(p ui.Element) {
		p.Margin(t.Space(1), 0, 0, 0).MaxWidth(t.Space(60)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border).Shadow(0, 6, 20, 0, kit.Shade(0.18))
		kit.Arrive(p, ui.Pose{Y: -t.Space(1), Scale: 0.96, Fade: 1}, 0)
		ui.Text(c, tip)
	})
	return root
}
