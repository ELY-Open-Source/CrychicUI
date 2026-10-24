package agent

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VoiceOrbOptions name the orb.
type VoiceOrbOptions struct {
	Label string
}

// VoiceOrb is an orb whose size follows a level.
func VoiceOrb(c *ui.Context, level float64, o VoiceOrbOptions) ui.Element {
	if !(level >= 0 && level <= 1) {
		panic(fmt.Sprintf("agent: VoiceOrb needs a level from 0 to 1, got %v", level))
	}
	def(&o.Label, "Voice level")
	t := c.Theme()
	side := t.Space(24)
	root := ui.Box(c).Size(side, side).Shrink(0).Role(ui.RoleMeter).Label(o.Label).Range(0, 1, level)
	at := root.AnimateTo("level", float32(level), ui.Bouncy)
	breath := float32(0.5 + 0.5*math.Sin(2*math.Pi*float64(root.Loop("breath", 2400*time.Millisecond, ui.Linear))))
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		cx, cy, full := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)/2
		core := full * (0.34 + 0.18*at)
		for i, a := range []float32{0.08, 0.14, 0.22} {
			rad := core + (full-core)*(1-float32(i)/3)*(0.45+0.55*at) + 2*breath
			p.Fill(ui.Rect{X: cx - rad, Y: cy - rad, W: 2 * rad, H: 2 * rad}, t.Accent.Alpha(a), 999)
		}
		p.Fill(ui.Rect{X: cx - core, Y: cy - core, W: 2 * core, H: 2 * core}, t.Accent, 999)
		hi := core * 0.45
		p.Fill(ui.Rect{X: cx - core*0.5, Y: cy - core*0.6, W: hi, H: hi}, t.AccentText.Alpha(0.35), 999)
	}))
	return root
}
