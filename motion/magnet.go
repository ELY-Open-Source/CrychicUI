package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Magnet pulls content toward the pointer while it is near.
func Magnet(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Magnet needs content")
	}
	root := ui.Box(c).Padding(c.Theme().Space(6))
	size := func() []float32 { b := root.LayoutBounds(); return []float32{b.W, b.H} }
	used := size()
	kit.Settle(c, root, used, size)
	x, y, over := root.PointerPosition()
	var dx, dy float32
	a := ui.Bouncy
	if over {
		dx, dy, a = (x-used[0]/2)*0.25, (y-used[1]/2)*0.25, ui.Snappy
	}
	return root.Children(func() {
		in := ui.Box(c)
		in.Translate(in.AnimateTo("x", dx, a), in.AnimateTo("y", dy, a)).Children(content)
	})
}
