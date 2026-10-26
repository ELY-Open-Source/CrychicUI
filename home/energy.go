package home

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EnergyOptions word the chart; Labels name each bar.
type EnergyOptions struct {
	Title, Unit, Peak string
	Labels            []string
}

// Energy charts energy use per period; hover reads a bar.
func Energy(c *ui.Context, use []float64, o EnergyOptions) ui.Element {
	def(&o.Title, "Energy today")
	def(&o.Unit, "kWh")
	def(&o.Peak, "Peak %s")
	if len(use) < 2 || len(use) > 48 || len(o.Labels) != len(use) {
		panic(fmt.Sprintf("home: Energy %q needs 2 to 48 values, each with a label", o.Title))
	}
	unique("Energy", o.Labels)
	peak, total := 0, 0.0
	for i, v := range use {
		if !(v >= 0 && v <= 1e6) {
			panic(fmt.Sprintf("home: Energy %q value %q must be 0 to 1e6", o.Title, o.Labels[i]))
		}
		if v > use[peak] {
			peak = i
		}
		total += v
	}
	top := math.Max(use[peak], 1e-9)
	t := c.Theme()
	kwh := func(v float64) string { return strconv.FormatFloat(v, 'f', 1, 64) }
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "zap", o.Title, func() {
			// A bounded chip: long labels truncate, the title keeps room.
			ui.Row(c).MinWidth(0).MaxWidth(t.Space(40)).Padding(t.Space(0.25), t.Space(2)).Radius(999).Background(t.Warning.Alpha(0.12)).
				Children(func() {
					ui.Text(c, fmt.Sprintf(o.Peak, o.Labels[peak])).FontSize(t.Rem(0.7)).Bold().TextColor(t.Warning).MinWidth(0).SingleLine()
				})
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Children(func() {
			hover := -1
			chart := ui.Row(c).FillWidth().Height(t.Space(26)).Gap(t.Space(0.75)).AlignItems(ui.Stretch)
			grow := chart.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Children(func() {
				for i, v := range use {
					c.Key(i)
					col := ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight()
					on := col.Hovered()
					if on {
						hover = i
					}
					h := col.AnimateTo("h", float32(v/top), ui.Smooth)
					base := t.Accent
					if i == peak {
						base = t.Warning
					}
					ink := col.AnimateColor("ink", base.Alpha(0.5+0.5*unit(on)), ui.Snappy)
					col.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						bh := max(t.Space(0.5), r.H*h*grow)
						p.Fill(ui.Rect{X: r.X, Y: r.Y + r.H - bh, W: r.W, H: bh}, ink, t.Space(0.5))
					}))
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
				for i, at := range [3]int{0, len(use) / 2, len(use) - 1} {
					align := [3]ui.Align{ui.Start, ui.Center, ui.End}[i]
					kit.Muted(c, o.Labels[at]).FontSize(t.Rem(0.62)).Grow(1).Basis(0).MinWidth(0).TextAlign(align).SingleLine()
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.End).Padding(t.Space(2), 0, 0, 0).Children(func() {
				value, note := total, o.Title
				if hover >= 0 {
					value, note = use[hover], o.Labels[hover]
				}
				kit.Numeric(c, kwh(value)).FontSize(t.Rem(1.4)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
				kit.Muted(c, o.Unit).FontSize(t.Rem(0.74)).Shrink(0)
				ui.Spacer(c)
				kit.Muted(c, note).FontSize(t.Rem(0.74)).MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root
}
