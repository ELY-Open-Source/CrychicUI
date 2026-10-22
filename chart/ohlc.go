package chart

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// OHLCOptions tune an OHLC.
type OHLCOptions = CandlesOptions

// OHLC ticks open and close on range lines.
func OHLC(c *ui.Context, label string, candles []Candle, at *int, o OHLCOptions) ui.Element {
	labels, lo, hi := ohlc("OHLC", label, candles)
	chosen("OHLC", label, at, len(candles))
	def(&o.Format, "%g")
	def(&o.Tip, "%s: open %s, high %s, low %s, close %s")
	t := c.Theme()
	x := mc.NewBand(labels...).Padding(0.35, 0.15)
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	xs := make([]float32, len(candles))
	for i := range candles {
		xs[i] = x.At(i)
	}
	describe := func(i int) string { return worded(o, candles[i]) }
	root := ui.Column(c).FillWidth()
	hover := ui.Local(root, "hover", func() int { return -1 })
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks()}, o.Height, func(area ui.Element) {
			g := float32(grow(area))
			shown := sprung(area, y, candles)
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				for i, k := range shown {
					ink := t.Success
					if candles[i].Close < candles[i].Open {
						ink = t.Danger
					}
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					mid, m := mc.X(r, x.At(i)), (k[0]+k[3])/2
					w := (mc.X(r, x.Band(i).Hi) - mc.X(r, x.Band(i).Lo)) / 2
					open, shut := mc.Y(r, m+(k[0]-m)*g), mc.Y(r, m+(k[3]-m)*g)
					p.Line(mid, mc.Y(r, m+(k[1]-m)*g), mid, mc.Y(r, m+(k[2]-m)*g), 2, ink)
					p.Line(mid-w, open, mid, open, 2, ink)
					p.Line(mid, shut, mid+w, shut, 2, ink)
				}
			}))
			follow(c, root, area, at, hover, xs, 1, describe, "OHLC")
		})
	})
}
