package finance

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NetPoint is one net-worth point in time.
type NetPoint struct {
	Label         string
	Assets, Debts float64
}

// NetWorthOptions word the card.
type NetWorthOptions struct{ Value, Assets, Debts string }

// NetWorth charts assets minus debts over time.
func NetWorth(c *ui.Context, points []NetPoint, currency string, o NetWorthOptions) ui.Element {
	if len(points) < 2 || currency == "" {
		panic("finance: NetWorth needs two points and a currency")
	}
	values := make([]float64, len(points))
	for i, p := range points {
		if p.Label == "" || p.Assets < 0 || p.Debts < 0 || !finite(p.Assets) || !finite(p.Debts) {
			panic(fmt.Sprintf("finance: NetWorth point %q needs assets and debts of zero or more", p.Label))
		}
		values[i] = p.Assets - p.Debts
	}
	lo, hi := slices.Min(values), slices.Max(values)
	span := hi - lo
	if !finite(span) {
		panic("finance: NetWorth range overflows")
	}
	def(&o.Value, "Net worth")
	def(&o.Assets, "Assets")
	def(&o.Debts, "Debts")
	first, last := values[0], values[len(values)-1]
	change := 0.0
	if first != 0 {
		change = 100 * ((last - first) / abs(first))
	}
	if !finite(change) {
		panic("finance: NetWorth change overflows")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Value).FontSize(t.Rem(0.8))
				money(c, last, currency).Bold().FontSize(t.Rem(1.8)).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			if first != 0 {
				changePill(c, change)
			}
		})
		ink := changeInk(t, last-first)
		if last == first {
			ink = t.Accent
		}
		chart := ui.Box(c).FillWidth().Height(t.Space(32))
		draw := chart.Keyframes("crychic.draw", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			for g := range 4 {
				p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H*float32(g)/3, W: r.W, H: 1}, t.Border.Alpha(0.6), 0)
			}
			pad := t.Space(1.5)
			at := func(i int) (float32, float32) {
				x := r.X + pad + (r.W-2*pad)*float32(i)/float32(len(values)-1)
				y := r.Y + r.H/2
				if hi > lo {
					y = r.Y + pad + (r.H-2*pad)*(1-float32((values[i]-lo)/span))
				}
				return x, y
			}
			var line, area ui.Path
			x0, _ := at(0)
			xn, _ := at(len(values) - 1)
			area.MoveTo(x0, r.Y+r.H)
			for i := range values {
				x, y := at(i)
				if i == 0 {
					line.MoveTo(x, y)
				} else {
					line.LineTo(x, y)
				}
				area.LineTo(x, y)
			}
			area.LineTo(xn, r.Y+r.H).Close()
			p.FillPath(&area, ink.Alpha(0.12*draw))
			p.StrokePath(line.Trim(0, draw), t.Space(0.5), ink)
			if draw >= 1 {
				x, y := at(len(values) - 1)
				d := t.Space(2.5)
				p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, ink, 999)
			}
		}))
		ui.Row(c).FillWidth().Children(func() {
			for i, p := range points {
				if len(points) > 7 && i != 0 && i != len(points)-1 {
					continue
				}
				c.Key(i)
				label := kit.Muted(c, p.Label).FontSize(t.Rem(0.7)).Grow(1).Basis(0).MinWidth(0).SingleLine()
				if i == len(points)-1 {
					label.TextAlign(ui.End)
				} else if i > 0 {
					label.TextAlign(ui.Center)
				}
			}
		})
		latest := points[len(points)-1]
		ui.Row(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(2.5), 0, 0, 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			for i, part := range []struct {
				name  string
				value float64
				ink   ui.Color
			}{{o.Assets, latest.Assets, t.Success}, {o.Debts, latest.Debts, t.Danger}} {
				c.Key(i)
				ui.Row(c).Grow(1).Basis(0).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(part.ink)
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						kit.Muted(c, part.name).FontSize(t.Rem(0.75))
						money(c, part.value, currency).Bold().FontSize(t.Rem(0.85))
					})
				})
			}
		})
	})
	return root
}
