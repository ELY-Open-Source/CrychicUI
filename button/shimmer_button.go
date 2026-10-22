package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShimmerButton sweeps a sheen across itself in a loop.
func ShimmerButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface).Clip()
	phase := root.Loop("shimmer", 1400*time.Millisecond, ui.Linear)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		w := r.W / 3
		x := r.X - w + (r.W+2*w)*phase
		g := ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{
			{Color: t.Accent.Alpha(0), At: 0},
			{Color: t.AccentText.Alpha(0.18), At: 0.5},
			{Color: t.Accent.Alpha(0), At: 1},
		}}
		p.FillGradient(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, g, 0)
	}))
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root.Children(func() { ui.Text(c, label).SingleLine() })
}
