package motion

import "github.com/egoist/mygo/ui"

// Presence scales and fades content in and out.
func Presence(c *ui.Context, shown bool, content func()) ui.Element {
	if content == nil {
		panic("motion: Presence needs content")
	}
	t := c.Theme()
	pop := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Scale: 0.95, Y: t.Space(2)},
		Exit: &ui.Motion{Scale: 0.95, Y: t.Space(2)}}
	return ui.Box(c).Children(func() {
		if shown {
			ui.Box(c).Transition(pop).Children(content)
		}
	})
}
