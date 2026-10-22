package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var fall = ui.Animation{Duration: 700 * time.Millisecond, Ease: ui.EaseIn}

// LevelMeter is audio-style level bars whose peak falls back.
func LevelMeter(c *ui.Context, label string, level float64) ui.Element {
	if label == "" {
		panic("feedback: LevelMeter needs a label")
	}
	v := fraction("LevelMeter", level)
	t := c.Theme()
	const n = 12
	root := ui.Row(c).Gap(t.Space(0.5)).Height(t.Space(5)).AlignItems(ui.End).AlignSelf(ui.Start).
		Role(ui.RoleMeter).Label(label).Range(0, 1, level)
	shown := root.AnimateTo("level", v, ui.Animation{Duration: 60 * time.Millisecond})
	held := root.AnimateTo("peak", v, fall)
	if v >= held {
		held = root.Jump("peak", v)
	}
	return root.Children(func() {
		for i := range n {
			c.Key(i)
			at := float32(i+1) / n
			tint := t.Success
			switch {
			case at > 0.9:
				tint = t.Danger
			case at > 0.7:
				tint = t.Warning
			}
			lit := shown >= at || held >= at && held < at+1.0/n
			bar := ui.Box(c).Width(t.Space(1.5)).Height(t.Space(1) + t.Space(4)*at).Radius(t.Space(0.25))
			if lit {
				bar.Background(tint)
			} else {
				bar.Background(tint.Alpha(0.15))
			}
		}
	})
}
