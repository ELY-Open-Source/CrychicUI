package motion

import "github.com/egoist/mygo/ui"

// ScaleIn grows content in from 0.9 on first build.
func ScaleIn(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: ScaleIn needs content")
	}
	root := ui.Box(c)
	k := root.Keyframes("crychic.in", arrive)
	s := 0.9 + 0.1*k
	return root.Scale(s, s).Opacity(k).Children(content)
}
