package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TiltCard tilts toward the pointer, then springs flat.
func TiltCard(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: TiltCard needs content")
	}
	t := c.Theme()
	root := kit.Card(c).Padding(t.Space(4))
	size := func() []float32 { b := root.LayoutBounds(); return []float32{b.W, b.H} }
	used := size()
	kit.Settle(c, root, used, size)
	x, y, over := root.PointerPosition()
	var rx, ry float32
	a := ui.Bouncy
	if over && used[0] > 0 && used[1] > 0 {
		ry = (x/used[0] - 0.5) * 16
		rx = -(y/used[1] - 0.5) * 16
		a = ui.Snappy
	}
	return root.RotateX(root.AnimateTo("rx", rx, a)).RotateY(root.AnimateTo("ry", ry, a)).Perspective(t.Space(200)).Children(content)
}
