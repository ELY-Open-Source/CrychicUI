package chart

import (
	"fmt"

	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// BulletOptions tune a Bullet.
type BulletOptions struct {
	// Bands are rising range tops behind the bar.
	Bands []float64
	// Format prints values; Target names the mark.
	Format, Target string
}

// Bullet sets a bar against a target and ranges.
func Bullet(c *ui.Context, label string, value, target float64, o BulletOptions) ui.Element {
	if label == "" {
		panic("chart: Bullet needs a label")
	}
	finite("Bullet", label, append([]float64{value, target}, o.Bands...)...)
	top := max(value, target, 0)
	for i, b := range o.Bands {
		if b <= 0 || (i > 0 && b <= o.Bands[i-1]) {
			panic(fmt.Sprintf("chart: Bullet %q bands must rise above zero", label))
		}
		top = max(top, b)
	}
	if value < 0 || target < 0 || top == 0 {
		panic(fmt.Sprintf("chart: Bullet %q values must be zero or more, one above", label))
	}
	def(&o.Format, "%g")
	def(&o.Target, "target")
	t := c.Theme()
	x := mc.NewLinear(0, top).Nice(4)
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		plot(c, label, mc.Axes{Bottom: x.Ticks(4)}, t.Space(14), func(area ui.Element) {
			g := float64(grow(area))
			lo := 0.0
			for i, b := range o.Bands {
				c.Key(fmt.Sprint("band", i))
				mc.Rect(c, x.Span(lo, b), mc.Span{Lo: 0, Hi: 1}).Background(t.TextMuted.Alpha(0.08 + 0.08*float32(len(o.Bands)-i)))
				lo = b
			}
			tip := fmt.Sprintf("%s: "+o.Format+", %s "+o.Format, label, value, o.Target, target)
			c.Key("bar")
			mc.Rect(c, x.Span(0, value*g), mc.Span{Lo: 0.3, Hi: 0.7}).Radius(t.Radius / 2).Background(t.Accent).
				Label(tip).Tooltip(tip).Transition(spring)
			c.Key("target")
			mc.Rect(c, x.Span(target, target), mc.Span{Lo: 0.15, Hi: 0.85}).MinWidth(3).Background(t.Text).Transition(spring)
		})
	})
}
