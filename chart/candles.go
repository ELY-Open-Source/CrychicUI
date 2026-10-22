package chart

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// Candle is a period's open, high, low and close.
type Candle struct {
	Label                  string
	Open, High, Low, Close float64
}

// CandlesOptions tune a Candles or an OHLC.
type CandlesOptions struct {
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
	// Tip words a candle: label, open, high, low, close.
	Tip string
}

// Candles draws bodies from open to close on wicks.
func Candles(c *ui.Context, label string, candles []Candle, at *int, o CandlesOptions) ui.Element {
	labels, lo, hi := ohlc("Candles", label, candles)
	chosen("Candles", label, at, len(candles))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: open %s, high %s, low %s, close %s")
	t := c.Theme()
	x := mc.NewBand(labels...).Padding(0.35, 0.15)
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float64(grow(area))
			shown := sprung(area, candles)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, k := range shown {
					mid, m := mc.X(r, x.At(i)), (k.Open+k.Close)/2
					p.Line(mid, mc.Y(r, y.At(m+(k.High-m)*g)), mid, mc.Y(r, y.At(m+(k.Low-m)*g)), 1, t.TextMuted)
				}
			}))
			for i, k := range candles {
				c.Key(i)
				ink := t.Success
				if k.Close < k.Open {
					ink = t.Danger
				}
				m := (k.Open + k.Close) / 2
				span := y.Span(m+(min(k.Open, k.Close)-m)*g, m+(max(k.Open, k.Close)-m)*g)
				span.Hi = max(span.Hi, span.Lo+0.004)
				mark(c, root, at, i, x.Band(i), span, ink, worded(o, k), "Candles")
			}
		})
	})
}
