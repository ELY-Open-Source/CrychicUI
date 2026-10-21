package choice

import (
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

var drawIn = ui.Animation{Duration: 160 * time.Millisecond}

// unit converts on to an animation target.
func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// box is a check box whose mark strokes in.
func box(c *ui.Context, row ui.Element, on, mixed bool) ui.Element {
	t := c.Theme()
	b := ui.Box(c).Size(t.Space(4), t.Space(4)).Radius(t.Space(1)).Shrink(0)
	fill, edge := t.Background, t.Border.Mix(t.Text, 0.25)
	if on || mixed {
		fill, edge = t.Accent, t.Accent
	} else if row.Hovered() {
		edge = t.Accent
	}
	b.Background(b.AnimateColor("fill", fill, ui.Snappy)).Border(1, b.AnimateColor("edge", edge, ui.Snappy))
	drawn := b.AnimateTo("mark", unit(on || mixed), drawIn)
	dash := b.AnimateTo("dash", unit(mixed), ui.Snappy)
	return b.DrawOver(func(p *ui.Painter, r ui.Rect) { mark(p, r, drawn, dash, t.Space(0.5), t.AccentText) })
}

// mark strokes a check, part drawn, morphed toward a dash.
func mark(p *ui.Painter, r ui.Rect, drawn, dash, width float32, col ui.Color) {
	if drawn <= 0 {
		return
	}
	check := [3][2]float32{{0.22, 0.52}, {0.42, 0.71}, {0.78, 0.31}}
	flat := [3][2]float32{{0.25, 0.5}, {0.5, 0.5}, {0.75, 0.5}}
	var pts [3][2]float32
	for i := range pts {
		pts[i] = [2]float32{
			r.X + r.W*(check[i][0]+(flat[i][0]-check[i][0])*dash),
			r.Y + r.H*(check[i][1]+(flat[i][1]-check[i][1])*dash),
		}
	}
	a, b := span(pts[0], pts[1]), span(pts[1], pts[2])
	left := (a + b) * min(drawn, 1)
	line := [][2]float32{pts[0]}
	if left <= a {
		line = append(line, along(pts[0], pts[1], left/a))
	} else {
		line = append(line, pts[1], along(pts[1], pts[2], (left-a)/b))
	}
	p.StrokePolyline(line, width, col)
}

func span(a, b [2]float32) float32 {
	return float32(math.Hypot(float64(b[0]-a[0]), float64(b[1]-a[1])))
}

func along(a, b [2]float32, f float32) [2]float32 {
	return [2]float32{a[0] + (b[0]-a[0])*f, a[1] + (b[1]-a[1])*f}
}

// dot is a radio circle whose center springs in.
func dot(c *ui.Context, row ui.Element, on bool) ui.Element {
	t := c.Theme()
	d := ui.Box(c).Size(t.Space(4), t.Space(4)).Radius(t.Space(2)).Shrink(0)
	fill, edge := t.Background, t.Border.Mix(t.Text, 0.25)
	if on {
		fill, edge = t.Accent, t.Accent
	} else if row.Hovered() {
		edge = t.Accent
	}
	d.Background(d.AnimateColor("fill", fill, ui.Snappy)).Border(1, d.AnimateColor("edge", edge, ui.Snappy))
	k := d.AnimateTo("dot", unit(on), ui.Bouncy)
	return d.DrawOver(func(p *ui.Painter, r ui.Rect) {
		s := t.Space(1.5) * max(k, 0)
		p.Fill(ui.Rect{X: r.X + (r.W-s)/2, Y: r.Y + (r.H-s)/2, W: s, H: s}, t.AccentText, s/2)
	})
}

// track is a switch track whose knob springs across.
func track(c *ui.Context, on bool) ui.Element {
	t := c.Theme()
	w, h, in := t.Space(9), t.Space(5), t.Space(0.5)
	tr := ui.Box(c).Size(w, h).Radius(h / 2).Shrink(0)
	off := t.Border.Mix(t.Text, 0.15)
	tr.Background(tr.AnimateColor("track", off.Mix(t.Accent, unit(on)), ui.Snappy))
	return tr.Children(func() {
		knob := ui.Box(c).Absolute().Top(in).Left(in).Size(h-2*in, h-2*in).Radius(h).Background(t.AccentText).
			Shadow(0, 1, 3, 0, kit.Shade(0.25))
		knob.Translate(knob.AnimateTo("x", (w-h)*unit(on), ui.Snappy), 0)
	})
}

// distinct panics unless items are non-empty and unique.
func distinct(component string, items []string) {
	if len(items) == 0 {
		panic("choice: " + component + " needs items")
	}
	for i, item := range items {
		if item == "" || slices.Contains(items[:i], item) {
			panic("choice: " + component + " needs distinct, non-empty items")
		}
	}
}

// subset panics unless selected holds distinct items.
func subset(component string, items []string, selected *[]string) {
	if selected == nil {
		panic("choice: " + component + " needs a selection")
	}
	for i, s := range *selected {
		if !slices.Contains(items, s) || slices.Contains((*selected)[:i], s) {
			panic("choice: " + component + " selection holds an unknown or repeated item")
		}
	}
}

// flip adds item when absent and removes it when present.
func flip(items []string, item string) []string {
	if i := slices.Index(items, item); i >= 0 {
		return slices.Delete(slices.Clone(items), i, i+1)
	}
	return append(slices.Clone(items), item)
}
