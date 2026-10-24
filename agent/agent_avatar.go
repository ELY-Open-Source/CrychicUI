package agent

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AgentAvatar is an avatar whose ring pulses while working.
func AgentAvatar(c *ui.Context, name string, image *ui.Bitmap, working bool) ui.Element {
	if name == "" {
		panic(fmt.Sprintf("agent: AgentAvatar %q needs a name", name))
	}
	t := c.Theme()
	w := t.Space(0.75)
	root := ui.Box(c).AlignSelf(ui.Start).Padding(t.Space(1.25)).Label(name).Role(ui.RoleImage)
	on := root.AnimateTo("on", unit(working), ui.Smooth)
	turn, glow := float32(0), float32(0)
	if working {
		turn = 360 * float32(root.Loop("turn", 1600*time.Millisecond, ui.Linear))
		glow = float32(0.5 + 0.5*math.Sin(2*math.Pi*float64(root.Loop("pulse", 1400*time.Millisecond, ui.Linear))))
	}
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)/2-w/2
		var ring, arc ui.Path
		p.StrokePath(ring.Circle(cx, cy, rad), w, t.Border.Mix(t.Accent.Alpha(0.3+0.4*glow), on))
		if on > 0 {
			p.Transform(ui.Rotation(turn), cx, cy, func() { p.StrokePath(arc.Circle(cx, cy, rad).Trim(0, 0.3*on), w, t.Accent) })
		}
	}))
	return root.Children(func() {
		display.Avatar(c, name, image)
		badge := ui.Box(c).Absolute().Right(0).Bottom(0).Size(t.Space(4.5), t.Space(4.5)).Center().Radius(999).
			Background(t.Accent).TextColor(t.AccentText).Border(2, t.Background)
		badge.Scale(on, on).Opacity(on).Children(func() { ui.Icon(c, icon.Must("sparkles")).Size(10, 10) })
	})
}
