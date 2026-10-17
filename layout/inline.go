package layout

import "github.com/egoist/mygo/ui"

// Inline is a wrapping row in theme gaps; size animates.
func Inline(c *ui.Context, gap float32, content func()) ui.Element {
	if !finite(gap) || gap < 0 {
		panic("layout: Inline gap must be finite and non-negative")
	}
	return ui.Row(c).Wrap().Gap(c.Theme().Space(gap)).Transition(ui.ElementTransition{Size: true}).Children(content)
}
