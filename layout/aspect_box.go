package layout

import "github.com/egoist/mygo/ui"

// AspectBox holds width over height at ratio; it resizes smoothly.
func AspectBox(c *ui.Context, ratio float32, content func()) ui.Element {
	if !finite(ratio) || ratio <= 0 {
		panic("layout: AspectBox ratio must be finite and positive")
	}
	return ui.Box(c).FillWidth().AspectRatio(ratio).Clip().Transition(ui.ElementTransition{Size: true}).
		Children(content)
}
