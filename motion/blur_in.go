package motion

import "github.com/egoist/mygo/ui"

// BlurIn sharpens content from a blur on first build.
func BlurIn(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: BlurIn needs content")
	}
	root := ui.Box(c)
	if k := root.Keyframes("crychic.in", arrive); k < 1 {
		root.Blur(c.Theme().Space(2) * (1 - k))
	}
	return root.Children(content)
}
