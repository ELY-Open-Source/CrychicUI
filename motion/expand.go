package motion

import "github.com/egoist/mygo/ui"

// Expand animates its height to fit content as content changes.
func Expand(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("motion: Expand needs content")
	}
	return ui.Column(c).ClipY().Transition(ui.ElementTransition{Size: true, Animation: &ui.Smooth}).Children(content)
}
