package motion

import "github.com/egoist/mygo/ui"

// HoverZoom scales content inside a clip while hovered.
func HoverZoom(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: HoverZoom needs content")
	}
	root := ui.Box(c).Clip().Radius(c.Theme().Radius)
	s := float32(1)
	if root.Hovered() {
		s = 1.06
	}
	return root.Children(func() {
		in := ui.Box(c)
		k := in.AnimateTo("zoom", s, ui.Smooth)
		in.Scale(k, k).Children(content)
	})
}
