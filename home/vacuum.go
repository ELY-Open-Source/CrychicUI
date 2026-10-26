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

// VacuumState is where a robot vacuum is in its run.
type VacuumState int

const (
	Docked VacuumState = iota
	Cleaning
	Returning
)

// VacuumOptions word the vacuum; States run docked to returning.
type VacuumOptions struct {
	Room, Start, Dock, Battery string
	States                     [3]string
}

// Vacuum starts and docks a robot vacuum; cue Tap.
func Vacuum(c *ui.Context, name string, state *VacuumState, battery float64, o VacuumOptions) ui.Element {
	if strings.TrimSpace(name) == "" || state == nil || *state < Docked || *state > Returning || !(battery >= 0 && battery <= 1) {
		panic(fmt.Sprintf("home: Vacuum %q needs a known state and a battery of 0 to 1", name))
	}
	def(&o.Start, "Start")
	def(&o.Dock, "Dock")
	def(&o.Battery, "%d%%")
	for i, s := range [...]string{"Docked", "Cleaning", "Returning"} {
		def(&o.States[i], s)
	}
	t := c.Theme()
	inks := [3]ui.Color{t.TextMuted, t.Accent, t.Warning}
	charge := t.Success
	switch {
	case battery < 0.2:
		charge = t.Danger
	case battery < 0.5:
		charge = t.Warning
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "bot", inks[*state], t.Space(10), *state != Docked)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
					c.Key(*state)
					kit.Arrive(ui.Text(c, o.States[*state]).FontSize(t.Rem(0.74)).Bold().TextColor(inks[*state]).Shrink(0), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
					if o.Room != "" {
						kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
					}
				})
			})
			ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(charge).Children(func() {
				glyph(c, "battery-charging", 16)
				kit.Numeric(c, fmt.Sprintf(o.Battery, int(math.Round(battery*100)))).FontSize(t.Rem(0.76)).Bold().Shrink(0)
			})
		})
		floor := ui.Box(c).FillWidth().Height(t.Space(24)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
		var sweep float32
		if *state == Cleaning {
			sweep = floor.Loop("sweep", 9*time.Second, ui.Linear)
		}
		home := floor.AnimateTo("home", unit(*state != Cleaning), ui.Animation{Duration: 1200 * time.Millisecond, Ease: ui.EaseInOut})
		floor.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			pad, rr := t.Space(3), t.Space(2.5)
			dock := ui.Rect{X: r.X + pad, Y: r.Y + r.H - pad - rr, W: rr * 2, H: rr}
			p.Fill(dock, t.Border, t.Space(0.5))
			// A lawnmower path: four passes across and back.
			lane := float64(sweep) * 4
			pass := math.Floor(lane)
			f := float32(lane - pass)
			if int(pass)%2 == 1 {
				f = 1 - f
			}
			x0, x1 := r.X+pad+rr, r.X+r.W-pad-rr
			y := r.Y + pad + rr + (r.H-2*pad-2*rr)*float32(pass)/3
			bx, by := x0+(x1-x0)*f, y
			dx, dy := dock.X+rr, dock.Y-rr*0.2
			x, yy := bx+(dx-bx)*home, by+(dy-by)*home
			p.FillPath(ring(x, yy, rr), t.Background)
			p.StrokePath(ring(x, yy, rr), t.Space(0.5), inks[*state])
			p.FillPath(ring(x, yy-rr*0.45, rr*0.22), inks[*state])
		}))
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			act := func(label, mark string, ok bool, next VacuumState) {
				c.Key(label)
				b := kit.Press(ui.ButtonBase(c).Label(label).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(2), t.Space(3)).Radius(999).
					Border(1, t.Border).Disabled(!ok))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						glyph(c, mark, 14)
						ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Shrink(0)
					})
				})
				if b.Clicked() && ok {
					*state = next
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("Vacuum", "state", "to", int(next))
				}
			}
			act(o.Start, "play", *state != Cleaning, Cleaning)
			act(o.Dock, "house", *state == Cleaning, Returning)
		})
	})
	return root
}
