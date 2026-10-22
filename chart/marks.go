package chart

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// mark is a focusable bar that tips, springs and picks.
func mark(c *ui.Context, root ui.Element, at *int, i int, x, y mc.Span, ink ui.Color, tip, comp string) ui.Element {
	t := c.Theme()
	b := mc.Rect(c, x, y).Radius(t.Radius / 2).Focusable().Label(tip).Tooltip(tip).Transition(spring)
	if *at >= 0 && *at != i {
		ink = ink.Alpha(0.4)
	}
	if b.Hovered() {
		ink = ink.Mix(t.Text, 0.15)
	}
	b.Background(b.AnimateColor("ink", ink, ui.Snappy))
	if b.Clicked() {
		pick(root, at, i, comp)
	}
	return b
}

// dots places points, sized by Size when sized.
func dots(c *ui.Context, label string, points []Point, at *int, o LineOptions, sized bool, comp string) ui.Element {
	var names []string
	var xs, ys []float64
	top := 0.0
	for _, pt := range points {
		finite(comp, label, pt.X, pt.Y, pt.Size)
		if pt.Size < 0 {
			panic(fmt.Sprintf("chart: %s %q sizes must not be negative", comp, label))
		}
		names, xs, ys = append(names, pt.Name), append(xs, pt.X), append(ys, pt.Y)
		top = max(top, pt.Size)
	}
	named(comp, label, names)
	chosen(comp, label, at, len(points))
	def(&o.Format, "%g")
	t := c.Theme()
	x := mc.NewLinear(extent(xs...)).Nice(5)
	y := mc.NewLinear(extent(ys...)).Nice(4)
	ticks := y.Ticks(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Left: ticks, Bottom: x.Ticks(5)}, o.Height, func(area ui.Element) {
			g := grow(area)
			area.Draw(grid(t, ticks))
			for i, pt := range points {
				c.Key(pt.Name)
				d := t.Space(3)
				if sized && top > 0 {
					d = t.Space(3) + t.Space(9)*float32(math.Sqrt(pt.Size/top))
				}
				tip := fmt.Sprintf("%s: "+o.Format+", "+o.Format, pt.Name, pt.X, pt.Y)
				ink := tone(t, i%4).Alpha(0.75)
				if *at >= 0 && *at != i {
					ink = ink.Alpha(0.3)
				}
				b := mc.Point(c, x.At(pt.X), y.At(pt.Y), d*g, d*g).Radius(999).Focusable().Label(tip).Tooltip(tip).Transition(spring)
				if b.Hovered() || *at == i {
					b.Border(2, t.Text.Alpha(0.6))
				}
				b.Background(b.AnimateColor("ink", ink, ui.Snappy))
				if b.Clicked() {
					pick(root, at, i, comp)
				}
			}
		})
	})
}

// sample checks groups and returns names and every value.
func sample(comp, label string, groups []Group) (names []string, all []float64) {
	for _, g := range groups {
		if len(g.Values) == 0 {
			panic(fmt.Sprintf("chart: %s %q groups need values", comp, label))
		}
		finite(comp, label, g.Values...)
		names, all = append(names, g.Name), append(all, g.Values...)
	}
	named(comp, label, names)
	return names, all
}

// ohlc checks candles and returns labels and the range.
func ohlc(comp, label string, candles []Candle) (labels []string, lo, hi float64) {
	if len(candles) == 0 {
		panic(fmt.Sprintf("chart: %s %q needs candles", comp, label))
	}
	lo, hi = candles[0].Low, candles[0].High
	for _, k := range candles {
		finite(comp, label, k.Open, k.High, k.Low, k.Close)
		if k.Low > min(k.Open, k.Close) || k.High < max(k.Open, k.Close) {
			panic(fmt.Sprintf("chart: %s %q candles need low and high around open and close", comp, label))
		}
		labels = append(labels, k.Label)
		lo, hi = min(lo, k.Low), max(hi, k.High)
	}
	named(comp, label, labels)
	if hi == lo {
		hi = lo + 1
	}
	return labels, lo, hi
}

// worded is a candle's tip in the caller's words.
func worded(o CandlesOptions, k Candle) string {
	f := func(v float64) string { return fmt.Sprintf(o.Format, v) }
	return fmt.Sprintf(o.Tip, k.Label, f(k.Open), f(k.High), f(k.Low), f(k.Close))
}

// sprung eases each candle's four values on e.
func sprung(e ui.Element, candles []Candle) []Candle {
	out := make([]Candle, len(candles))
	for i, k := range candles {
		v := springs(e, fmt.Sprint("k", i), []float64{k.Open, k.High, k.Low, k.Close})
		out[i] = Candle{k.Label, v[0], v[1], v[2], v[3]}
	}
	return out
}

// spark is a tiny drawing whose hovered value tips.
func spark(c *ui.Context, label string, vs []float64, o SparklineOptions, comp string, draw func(p *ui.Painter, r ui.Rect, ys []float32, g float32, hot int, ink ui.Color)) ui.Element {
	if label == "" || len(vs) == 0 {
		panic(fmt.Sprintf("chart: %s %q needs values", comp, label))
	}
	finite(comp, label, vs...)
	def(&o.Format, "%g")
	t := c.Theme()
	for _, v := range []float32{o.Width, o.Height} {
		if v < 0 || math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
			panic(fmt.Sprintf("chart: %s %q sizes must be finite and not negative", comp, label))
		}
	}
	if o.Width == 0 {
		o.Width = t.Space(30)
	}
	if o.Height == 0 {
		o.Height = t.Space(8)
	}
	lo, hi := mc.Extent(vs...)
	if hi == lo {
		hi = lo + 1
	}
	box := ui.Box(c).Size(o.Width, o.Height).Role(ui.RoleImage).Label(label)
	hot := -1
	if x, _, over := box.PointerPosition(); over {
		hot = min(int(x/o.Width*float32(len(vs))), len(vs)-1)
		box.Tooltip(fmt.Sprintf(o.Format, vs[hot]))
	}
	shown := springs(box, "v", vs)
	ys := make([]float32, len(vs))
	for i, v := range shown {
		ys[i] = float32((v - lo) / (hi - lo))
	}
	g := grow(box)
	return box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) { draw(p, r, ys, g, hot, t.Accent) }))
}

// values checks a value per label.
func values(comp, label string, labels []string, vs []float64) {
	named(comp, label, labels)
	if len(vs) != len(labels) {
		panic(fmt.Sprintf("chart: %s %q needs a value per label", comp, label))
	}
	finite(comp, label, vs...)
}

// table checks series of a value per label.
func table(comp, label string, labels []string, series []Series) (vals [][]float64, names []string) {
	named(comp, label, labels)
	if len(series) == 0 {
		panic(fmt.Sprintf("chart: %s %q needs series", comp, label))
	}
	for _, s := range series {
		if len(s.Values) != len(labels) {
			panic(fmt.Sprintf("chart: %s %q series need a value per label", comp, label))
		}
		finite(comp, label, s.Values...)
		vals, names = append(vals, s.Values), append(names, s.Name)
	}
	named(comp, label, names)
	return vals, names
}

// extent is lo to hi with zero, never empty.
func extent(vs ...float64) (float64, float64) {
	lo, hi := mc.Extent(vs...)
	lo, hi = min(lo, 0), max(hi, 0)
	if lo == hi {
		hi = lo + 1
	}
	return lo, hi
}
