package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VoiceOrb is an orb whose size follows a level.
func VoiceOrb(c *ui.Context, level float64) ui.Element {
	if !(level >= 0 && level <= 1) {
		panic(fmt.Sprintf("agent: VoiceOrb needs a level from 0 to 1, got %v", level))
	}
	t := c.Theme()
	side := t.Space(12)
	root := ui.Box(c).Size(side, side).Role(ui.RoleMeter).Label("Voice level").Range(0, 1, level)
	at := root.AnimateTo("level", float32(level), ui.Smooth)
	rad := side/2 - t.Space(0.5)
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		size := rad * (0.45 + 0.55*at)
		p.Fill(ui.Rect{X: r.X + r.W/2 - size, Y: r.Y + r.H/2 - size, W: 2 * size, H: 2 * size}, t.Accent, 999)
	}))
	return root
}
