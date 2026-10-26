package home

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FanOptions word the fan; Speeds run from off to high.
type FanOptions struct {
	Room   string
	Speeds [4]string
}

// Fan sets speed; its blades spin to match; cue Tick.
func Fan(c *ui.Context, name string, speed *int, o FanOptions) ui.Element {
	if strings.TrimSpace(name) == "" || speed == nil || *speed < 0 || *speed > 3 {
		panic(fmt.Sprintf("home: Fan %q needs a speed of 0 to 3", name))
	}
	for i, s := range [...]string{"Off", "Low", "Medium", "High"} {
		def(&o.Speeds[i], s)
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			size := t.Space(16)
			blades := ui.Box(c).Size(size, size).Shrink(0).Radius(999)
			blades.Background(blades.AnimateColor("bg", t.SurfaceHover.Mix(t.Accent, 0.1*unit(*speed > 0)), ui.Snappy))
			var turn float32
			if *speed > 0 {
				// Faster speeds turn in less time.
				turn = blades.Loop("spin", time.Duration(1400 / *speed)*time.Millisecond, ui.Linear)
			}
			ink := blades.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(*speed > 0)), ui.Snappy)
			blades.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2
				for i := 0; i < 3; i++ {
					a := float64(turn)*2*math.Pi + float64(i)*2*math.Pi/3
					bx, by := cx+rr*0.42*float32(math.Sin(a)), cy-rr*0.42*float32(math.Cos(a))
					p.FillPath(ring(bx, by, rr*0.3), ink.Alpha(0.85))
				}
				p.FillPath(ring(cx, cy, rr*0.16), ink)
			}))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					c.Key(*speed)
					kit.Arrive(ui.Text(c, o.Speeds[*speed]).FontSize(t.Rem(0.74)).Bold().TextColor(t.TextMuted.Mix(t.Accent, unit(*speed > 0))).Shrink(0),
						ui.Pose{Y: t.Space(1), Fade: 1}, 0)
					if o.Room != "" {
						kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
					}
				})
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(1)).Radius(999).Background(t.SurfaceHover).Children(func() {
			for i, s := range o.Speeds {
				on := i == *speed
				b := kit.Press(ui.ButtonBase(c).Label(s).Checked(on).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(1.5), t.Space(1)).Radius(999))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Mix(t.Background, unit(on || b.Hovered())), ui.Snappy))
				b.Children(func() {
					ui.Text(c, s).FontSize(t.Rem(0.74)).Bold().TextColor(t.TextMuted.Mix(t.Text, unit(on))).FillWidth().MinWidth(0).
						TextAlign(ui.Center).SingleLine()
				})
				if b.Clicked() && !on {
					*speed = i
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("Fan", "speed", "level", i)
				}
			}
		})
	})
	return root
}
