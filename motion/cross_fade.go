package motion

import "github.com/egoist/mygo/ui"

// CrossFade shows a or b, fading one into the other.
func CrossFade(c *ui.Context, showB bool, a, b func()) ui.Element {
	if a == nil || b == nil {
		panic("motion: CrossFade needs two contents")
	}
	fade := ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return ui.Box(c).Transition(ui.ElementTransition{Size: true}).Children(func() {
		c.Key(showB)
		w := ui.Box(c).Transition(fade)
		if showB {
			w.Children(b)
		} else {
			w.Children(a)
		}
	})
}
