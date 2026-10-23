package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TunerOptions tune a Tuner.
type TunerOptions struct {
	// Step is one detent; zero gives 0.1.
	Step float64
	// Format writes the frequency; "%.1f MHz" when empty.
	Format string
}

// Tuner is a dial over a sliding band.
func Tuner(c *ui.Context, label string, freq *float64, lo, hi float64, o TunerOptions) ui.Element {
	if o.Step == 0 {
		o.Step = 0.1
	}
	def(&o.Format, "%.1f MHz")
	t := c.Theme()
	places := max(0, int(math.Ceil(-math.Log10(o.Step)-1e-9)))
	root := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center)
	return root.Children(func() {
		if choice.Knob(c, label, freq, lo, hi, choice.KnobOptions{Step: o.Step, Format: fmt.Sprintf("%%.%df", places)}).Changed() {
			root.SetChanged()
		}
		ui.Column(c).Grow(1).Gap(t.Space(2)).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Format, *freq)).FontSize(t.Rem(1.6)).Bold().ContentTransition(ui.ContentTransition{})
			band := ui.Box(c).FillWidth().Height(t.Space(12)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
			at := band.AnimateTo("at", float32((*freq-lo)/o.Step), ui.Snappy)
			steps := int(math.Round((hi - lo) / o.Step))
			gap, needle := t.Space(1.5), t.Space(0.5)
			band.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
				mid := r.X + r.W/2
				for k := max(0, int(math.Floor(float64(at-r.W/2/gap)))); k <= steps; k++ {
					x := mid + (float32(k)-at)*gap
					if x > r.X+r.W {
						break
					}
					tall, v := r.H*0.25, lo+float64(k)*o.Step
					if int(math.Round(v/o.Step))%10 == 0 {
						tall = r.H * 0.45
						p.Text(x+t.Space(0.75), r.Y+t.Space(1), fmt.Sprintf("%.*f", places, v), t.Rem(0.65), t.TextMuted)
					}
					p.Fill(ui.Rect{X: x - 0.5, Y: r.Y + r.H - tall, W: 1, H: tall}, t.TextMuted, 0)
				}
				p.Fill(ui.Rect{X: mid - needle/2, Y: r.Y, W: needle, H: r.H}, t.Accent, needle/2)
			}))
		})
	})
}
