package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ParallaxOptions tune a Parallax.
type ParallaxOptions struct {
	// Height is the scene's height; t.Space(40) when zero.
	Height float32
}

// Parallax moves each layer against the pointer by its depth.
func Parallax(c *ui.Context, depths []float32, layer func(i int), o ParallaxOptions) ui.Element {
	if len(depths) == 0 || layer == nil {
		panic("motion: Parallax needs depths and a layer")
	}
	for _, d := range depths {
		finite("Parallax", float64(d))
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(positive("Parallax", "Height", o.Height, t.Space(40))).Clip().Radius(t.Radius)
	size := func() []float32 { b := root.LayoutBounds(); return []float32{b.W, b.H} }
	used := size()
	kit.Settle(c, root, used, size)
	x, y, over := root.PointerPosition()
	var dx, dy float32
	if over && used[0] > 0 && used[1] > 0 {
		dx, dy = x/used[0]*2-1, y/used[1]*2-1
	}
	shift := t.Space(4)
	return root.Children(func() {
		for i, d := range depths {
			c.Key(i)
			l := ui.Box(c).Absolute().Top(0).Bottom(0).Left(0).Right(0)
			l.Translate(l.AnimateTo("x", -dx*d*shift, ui.Smooth), l.AnimateTo("y", -dy*d*shift, ui.Smooth)).Children(func() { layer(i) })
		}
	})
}
