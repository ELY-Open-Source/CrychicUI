package chart

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// StreamingOptions tune a Streaming.
type StreamingOptions struct {
	// Window is how many values show; zero gives 30.
	Window int
	// Height of the chart; zero gives the default.
	Height float32
	// Format prints a value; "%g" when empty.
	Format string
}

// arrival is when Streaming last took new values.
type arrival struct {
	n, k int
	at   time.Time
}

// Streaming scrolls its line left as values arrive.
func Streaming(c *ui.Context, label string, vs []float64, o StreamingOptions) ui.Element {
	if label == "" || len(vs) == 0 {
		panic("chart: Streaming needs a label and values")
	}
	finite("Streaming", label, vs...)
	if o.Window < 0 || o.Window == 1 {
		panic(fmt.Sprintf("chart: Streaming %q window must be 2 or more", label))
	}
	if o.Window == 0 {
		o.Window = 30
	}
	def(&o.Format, "%g")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	st := ui.Local(root, "arrival", func() arrival { return arrival{n: len(vs)} })
	if len(vs) != st.n {
		st.k = max(0, min(len(vs)-st.n, o.Window))
		st.n, st.at = len(vs), c.Now()
		kit.Log("Streaming", "arrived", "count", st.k)
	}
	const dur = 300 * time.Millisecond
	shift := float32(0)
	if u := float32(c.Now().Sub(st.at)) / float32(dur); st.k > 0 && u < 1 && !c.Preferences().ReduceMotion {
		shift = float32(st.k) * (1 - ui.EaseOut(u))
		c.AnimationFrame()
	}
	seen := vs[max(0, len(vs)-o.Window-st.k):]
	lo, hi := mc.Extent(seen...)
	if hi == lo {
		hi = lo + 1
	}
	y := mc.NewLinear(lo, hi).Nice(4)
	ticks := y.Ticks(4)
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			kit.Muted(c, label).FontSize(t.Rem(0.8))
			kit.Numeric(c, fmt.Sprintf(o.Format, vs[len(vs)-1])).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
		})
		plot(c, label, mc.Axes{Left: ticks}, o.Height, func(area ui.Element) {
			g := grow(area)
			bot, top := area.AnimateTo("lo", float32(y.Lo), ui.Snappy), area.AnimateTo("hi", float32(y.Hi), ui.Snappy)
			xs, ys := make([]float32, len(seen)), make([]float32, len(seen))
			for j, v := range seen {
				xs[j] = (float32(j-(len(seen)-1)+o.Window-1) + shift) / float32(o.Window-1)
				ys[j] = (float32(v) - bot) / max(top-bot, 1e-6)
			}
			area.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				grid(t, ticks)(p, r)
				p.StrokePath(mc.Line(r, xs, ys, mc.CurveMonotone).Trim(0, g), 2, t.Accent)
				last := len(xs) - 1
				p.FillPath(new(ui.Path).Circle(mc.X(r, xs[last]), mc.Y(r, ys[last]), t.Space(1.25)), t.Accent)
			}))
		})
	})
}
