package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EqualizerOptions tune an Equalizer.
type EqualizerOptions struct {
	// Range bounds gains in dB; 12 by default.
	Range float64
	// Step snaps a gain; zero gives 1.
	Step float64
	// Height of a slider; zero gives the default.
	Height float32
}

// Equalizer edits band gains; the curve follows.
func Equalizer(c *ui.Context, label string, bands []string, gains []float64, o EqualizerOptions) ui.Element {
	if label == "" || len(bands) == 0 || len(gains) != len(bands) {
		panic(fmt.Sprintf("media: Equalizer %q needs a label and one gain per band", label))
	}
	if o.Range == 0 {
		o.Range = 12
	}
	if o.Step == 0 {
		o.Step = 1
	}
	if !(o.Range > 0) || math.IsInf(o.Range, 0) || !(o.Step > 0) || math.IsInf(o.Step, 0) {
		panic(fmt.Sprintf("media: Equalizer %q range and step must be finite and above 0", label))
	}
	for _, b := range bands {
		if b == "" {
			panic(fmt.Sprintf("media: Equalizer %q bands need names", label))
		}
	}
	t := c.Theme()
	size("Equalizer", label, &o.Height, t.Space(32))
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Label(label)
	return root.Children(func() {
		curve := ui.Box(c).FillWidth().Height(t.Space(16)).Radius(t.Radius).Background(t.Surface)
		gs := make([]float32, len(gains))
		for i, g := range gains {
			gs[i] = curve.AnimateTo(fmt.Sprint("gain.", i), float32(g/o.Range), ui.Snappy)
		}
		curve.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			mid := r.Y + r.H/2
			p.Fill(ui.Rect{X: r.X, Y: mid, W: r.W, H: 1}, t.Border, 0)
			at := func(i int) int { return max(0, min(len(gs)-1, i)) }
			x := func(i int) float32 { return r.X + r.W*(float32(at(i))+0.5)/float32(len(gs)) }
			y := func(i int) float32 { return mid - gs[at(i)]*r.H*0.42 }
			var line, area ui.Path
			line.MoveTo(r.X, y(0)).LineTo(x(0), y(0))
			area.MoveTo(r.X, mid).LineTo(r.X, y(0)).LineTo(x(0), y(0))
			for i := 0; i+1 < len(gs); i++ {
				c1x, c1y := x(i)+(x(i+1)-x(i-1))/6, y(i)+(y(i+1)-y(i-1))/6
				c2x, c2y := x(i+1)-(x(i+2)-x(i))/6, y(i+1)-(y(i+2)-y(i))/6
				line.CubeTo(c1x, c1y, c2x, c2y, x(i+1), y(i+1))
				area.CubeTo(c1x, c1y, c2x, c2y, x(i+1), y(i+1))
			}
			end := len(gs) - 1
			line.LineTo(r.X+r.W, y(end))
			area.LineTo(r.X+r.W, y(end)).LineTo(r.X+r.W, mid).Close()
			p.FillPath(&area, t.Accent.Alpha(0.15))
			p.StrokePath(&line, t.Space(0.5), t.Accent)
			d := t.Space(1.5)
			for i := range gs {
				p.Fill(ui.Rect{X: x(i) - d/2, Y: y(i) - d/2, W: d, H: d}, t.Accent, d/2)
			}
		}))
		ui.Row(c).FillWidth().Children(func() {
			for i := range bands {
				c.Key(i)
				ui.Column(c).Grow(1).Basis(0).AlignItems(ui.Center).Gap(t.Space(1)).Children(func() {
					s := choice.VerticalSlider(c, bands[i], &gains[i], -o.Range, o.Range,
						choice.VerticalSliderOptions{Height: o.Height, Step: o.Step})
					if s.Changed() {
						root.SetChanged()
						kit.Log("Equalizer", "changed", "band", i)
					}
					kit.Muted(c, bands[i]).FontSize(t.Rem(0.75)).SingleLine()
				})
			}
		})
	})
}
