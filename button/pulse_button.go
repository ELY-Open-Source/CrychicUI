package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PulseButton pulses a ring until clicked.
func PulseButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	clicked := ui.Local(root, "clicked", func() bool { return false })
	if root.Clicked() {
		*clicked = true
		sound.Play(sound.Tap)
	}
	if !*clicked {
		phase := root.Loop("pulse", time.Second, ui.Linear)
		pulse := phase * 2
		if pulse > 1 {
			pulse = 2 - pulse
		}
		root.DrawOver(func(p *ui.Painter, r ui.Rect) {
			grow := 2 + 2*pulse
			p.Stroke(ui.Rect{X: r.X - grow, Y: r.Y - grow, W: r.W + 2*grow, H: r.H + 2*grow},
				t.Accent.Alpha(0.2+0.5*pulse), t.Radius+grow, 1)
		})
	}
	return root.Children(func() { ui.Text(c, label).SingleLine() })
}
