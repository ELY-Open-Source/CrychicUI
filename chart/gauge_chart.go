package chart

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// GaugeChartOptions tune a GaugeChart.
type GaugeChartOptions struct {
	// Bands are rising band tops: green, amber, red.
	Bands []float64
	Size  float32
	// Format prints the value; "%g" when empty.
	Format string
}

// GaugeChart is a half arc with bands and a needle.
func GaugeChart(c *ui.Context, label string, value, lo, hi float64, o GaugeChartOptions) ui.Element {
	if label == "" {
		panic("chart: GaugeChart needs a label")
	}
	finite("GaugeChart", label, append([]float64{value, lo, hi}, o.Bands...)...)
	if hi <= lo || value < lo || value > hi {
		panic(fmt.Sprintf("chart: GaugeChart %q needs lo < hi holding the value", label))
	}
	for i, b := range o.Bands {
		if b <= lo || b > hi || (i > 0 && b <= o.Bands[i-1]) {
			panic(fmt.Sprintf("chart: GaugeChart %q bands must rise within lo to hi", label))
		}
	}
	def(&o.Format, "%g")
	t := c.Theme()
	size := o.Size
	if size == 0 {
		size = t.Space(56)
	}
	inks := []ui.Color{t.Success, t.Warning, t.Danger}
	tip := fmt.Sprintf("%s: "+o.Format, label, value)
	root := ui.Column(c).AlignItems(ui.Center).Gap(t.Space(1))
	return root.Children(func() {
		box := ui.Box(c).Size(size, size/2+t.Space(2)).Role(ui.RoleMeter).Label(label).Range(lo, hi, value).Tooltip(tip)
		frac := box.AnimateTo("needle", float32((value-lo)/(hi-lo)), ui.Bouncy) * grow(box)
		box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+size/2
			outer := size/2 - t.Space(1)
			arc := func(a, b float64, ink ui.Color) {
				s := mc.Slice{Start: -math.Pi/2 + math.Pi*(a-lo)/(hi-lo), End: -math.Pi/2 + math.Pi*(b-lo)/(hi-lo)}
				p.FillPath(mc.Sector(cx, cy, outer*0.78, outer, s), ink)
			}
			from := lo
			for i, b := range o.Bands {
				arc(from, b, inks[min(i, 2)].Alpha(0.75))
				from = b
			}
			rest := t.Border
			if len(o.Bands) > 0 {
				rest = inks[min(len(o.Bands), 2)].Alpha(0.75)
			}
			arc(from, hi, rest)
			a := -math.Pi/2 + math.Pi*float64(frac)
			tx, ty := mc.Polar(cx, cy, outer*0.9, a)
			var needle ui.Path
			needle.MoveTo(cx, cy).LineTo(tx, ty)
			p.StrokePath(&needle, 3, t.Text)
			p.FillPath(new(ui.Path).Circle(cx, cy, t.Space(1.5)), t.Text)
		}))
		kit.Numeric(c, fmt.Sprintf(o.Format, value)).FontSize(t.Rem(1.3)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
	})
}
