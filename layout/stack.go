package layout

import "github.com/egoist/mygo/ui"

// Stack is a column gapped in theme units; size animates.
func Stack(c *ui.Context, gap float32, content func()) ui.Element {
	if !finite(gap) || gap < 0 {
		panic("layout: Stack gap must be finite and non-negative")
	}
	return ui.Column(c).Gap(c.Theme().Space(gap)).Transition(ui.ElementTransition{Size: true}).Children(content)
}
