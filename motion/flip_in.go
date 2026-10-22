package motion

import "github.com/egoist/mygo/ui"

// FlipIn turns content in about its top edge.
func FlipIn(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: FlipIn needs content")
	}
	root := ui.Box(c).TransformOrigin(0.5, 0)
	k := root.Keyframes("crychic.in", arrive)
	return root.RotateX(-90 * (1 - k)).Perspective(c.Theme().Space(200)).Opacity(k).Children(content)
}
