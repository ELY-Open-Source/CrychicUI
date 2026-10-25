package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HealthBarOptions word the bar; Count takes hp and max.
type HealthBarOptions struct{ Label, Count string }

// HealthBar shows hit points; damage leaves a trail that drains.
func HealthBar(c *ui.Context, hp, most int, o HealthBarOptions) ui.Element {
	def(&o.Label, "HP")
	if most <= 0 || hp < 0 || hp > most {
		panic(fmt.Sprintf("game: HealthBar %q needs 0 <= hp <= max and a positive max", o.Label))
	}
	def(&o.Count, "%d / %d")
	t := c.Theme()
	share := float32(float64(hp) / float64(most))
	ink := t.Success
	switch {
	case share <= 0.25:
		ink = t.Danger
	case share <= 0.5:
		ink = t.Warning
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleMeter).Label(o.Label)
	prev := ui.Local(root, "prev", func() float32 { return share })
	held := ui.Local(root, "held", func() float32 { return share })
	hit := ui.Local(root, "hit", func() time.Time { return time.Time{} })
	now := c.Now()
	if share < *prev {
		// The trail holds where the bar was, then drains.
		*held = max(*held, *prev)
		if now.Sub(*hit) > 450*time.Millisecond {
			*held = *prev
		}
		*hit = now
	}
	*prev = share
	trail := share
	if wait := 450*time.Millisecond - now.Sub(*hit); wait > 0 {
		trail = *held
		c.After(wait)
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			g := glyph(c, "heart", 14)
			g.TextColor(g.AnimateColor("ink", ink, ui.Snappy))
			if share <= 0.25 && hp > 0 {
				s := 1 + 0.15*float32(g.Loop("beat", 900*time.Millisecond, ui.EaseInOut))
				g.Scale(s, s)
			}
			ui.Text(c, o.Label).FontSize(t.Rem(0.75)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Count, hp, most)).Font("monospace").FontSize(t.Rem(0.75)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
		bar := ui.Box(c).FillWidth().Height(t.Space(3.5)).Radius(999).Background(t.SurfaceHover)
		fill := bar.AnimateTo("fill", share, ui.Snappy)
		lag := bar.AnimateTo("trail", trail, ui.Smooth)
		col := bar.AnimateColor("ink", ink, ui.Snappy)
		bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			if lag > fill {
				p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * lag, H: r.H}, t.Danger.Alpha(0.35), 999)
			}
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, col, 999)
			p.Fill(ui.Rect{X: r.X + r.H/2, Y: r.Y + r.H*0.18, W: max(0, r.W*fill-r.H), H: r.H * 0.22}, t.Background.Alpha(0.3), 999)
		}))
	})
	return root
}
