package agent

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AgentAvatar is an avatar whose ring pulses while working.
func AgentAvatar(c *ui.Context, name string, image *ui.Bitmap, working bool) ui.Element {
	if name == "" {
		panic("agent: AgentAvatar needs a name")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(1)).Label(name)
	glow := float32(0.4)
	if working {
		ph := float64(root.Loop("pulse", 1400*time.Millisecond, ui.Linear))
		glow = 0.4 + 0.6*float32(math.Abs(math.Sin(math.Pi*ph)))
	}
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		p.Stroke(r, t.Accent.Alpha(glow), 999, t.Space(0.75))
	}))
	return root.Children(func() {
		display.Avatar(c, name, image)
	})
}
