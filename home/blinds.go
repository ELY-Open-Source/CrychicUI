package home

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BlindsOptions word the blinds; Level takes the percent open.
type BlindsOptions struct{ Room, Open, Closed, Level, Position string }

// Blinds slides blinds open and shut; cue Tick.
func Blinds(c *ui.Context, name string, open *float64, o BlindsOptions) ui.Element {
	if strings.TrimSpace(name) == "" || open == nil || !(*open >= 0 && *open <= 1) {
		panic(fmt.Sprintf("home: Blinds %q needs an open share of 0 to 1", name))
	}
	def(&o.Open, "Open")
	def(&o.Closed, "Closed")
	def(&o.Level, "%d%% open")
	def(&o.Position, "Position")
	t := c.Theme()
	pct := int(math.Round(*open * 100))
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3.5)).AlignItems(ui.Center).Children(func() {
			win := ui.Box(c).Size(t.Space(18), t.Space(22)).Shrink(0).Radius(t.Radius).Border(t.Space(0.75), t.Border)
			k := win.AnimateTo("open", float32(*open), ui.Smooth)
			win.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
				p.Fill(r, t.Accent.Alpha(0.12), t.Radius)
				p.Line(r.X+r.W/2, r.Y, r.X+r.W/2, r.Y+r.H, 1, t.Border)
				p.Line(r.X, r.Y+r.H/2, r.X+r.W, r.Y+r.H/2, 1, t.Border)
				const slats = 10
				drop := r.H * (1 - k)
				gap := drop / slats
				for i := 0; i < slats; i++ {
					y := r.Y + float32(i)*gap
					p.Fill(ui.Rect{X: r.X, Y: y, W: r.W, H: max(gap*0.8, t.Space(0.5))}, t.SurfaceHover.Mix(t.Border, 0.6), 0)
				}
				p.Fill(ui.Rect{X: r.X, Y: r.Y + drop, W: r.W, H: t.Space(1)}, t.TextMuted, 0)
			}))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
				if o.Room != "" {
					kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				}
				state := fmt.Sprintf(o.Level, pct)
				switch pct {
				case 0:
					state = o.Closed
				case 100:
					state = o.Open
				}
				kit.Numeric(c, state).FontSize(t.Rem(1)).Bold().Shrink(0).Padding(t.Space(1.5), 0, 0, 0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
		slider := ui.Slider(c, open, 0, 1).Step(0.05).Label(o.Position).FillWidth()
		if slider.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Blinds", "moved", "percent", int(math.Round(*open*100)))
		}
	})
	return root
}
