package devtools

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CPUOptions word the card; Core and Average are formats.
type CPUOptions struct{ Title, Core, Average string }

// CPU shows each core's load as a bar.
func CPU(c *ui.Context, cores []float64, o CPUOptions) ui.Element {
	if len(cores) == 0 {
		panic("devtools: CPU needs cores")
	}
	var total float64
	for i, v := range cores {
		if !finite(v) || v < 0 || v > 1 {
			panic(fmt.Sprintf("devtools: CPU core %d load %v is outside 0..1", i, v))
		}
		total += v
	}
	def(&o.Title, "CPU")
	def(&o.Core, "%d")
	def(&o.Average, "%d%% average")
	t := c.Theme()
	pct := func(v float64) int { return int(math.Round(v * 100)) }
	ink := func(v float64) ui.Color {
		switch {
		case v >= 0.9:
			return t.Danger
		case v >= 0.7:
			return t.Warning
		}
		return t.Accent
	}
	avg := total / float64(len(cores))
	cell, gap := t.Space(8), t.Space(1.5)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "cpu", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Average, pct(avg)), ink(avg)) })
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		view.Children(func() {
			row := ui.Row(c).Width(max(seen, float32(len(cores))*(cell+gap)-gap+t.Space(8))).Gap(gap).Padding(t.Space(4))
			row.Children(func() {
				for i, v := range cores {
					col := ui.Column(c).Grow(1).Basis(0).MinWidth(cell).Gap(t.Space(1)).AlignItems(ui.Center)
					kit.Stagger(c, row, col, i).Children(func() {
						kit.Numeric(c, fmt.Sprintf("%d%%", pct(v))).FontSize(t.Rem(0.62)).TextColor(t.TextMuted).Shrink(0)
						bar := ui.Box(c).FillWidth().Height(t.Space(24)).Radius(t.Radius * 0.6).Background(t.SurfaceHover)
						h := bar.AnimateTo("load", float32(v), ui.Smooth)
						tone := bar.AnimateColor("ink", ink(v), ui.Snappy)
						bar.Draw(kit.Inside(t.Radius*0.6, func(p *ui.Painter, r ui.Rect) {
							p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H*(1-h), W: r.W, H: r.H * h}, tone, t.Radius*0.6)
						}))
						kit.Muted(c, fmt.Sprintf(o.Core, i+1)).Font("monospace").FontSize(t.Rem(0.62)).Shrink(0)
					})
				}
			})
		})
	})
	return root
}
