package layout

import "github.com/egoist/mygo/ui"

// Collapse shows content while open; its height springs.
func Collapse(c *ui.Context, open bool, content func()) ui.Element {
	return ui.Column(c).FillWidth().Transition(ui.ElementTransition{Size: true, Animation: &collapseSpring}).Children(func() {
		if open {
			ui.Column(c).Key("body").FillWidth().Transition(collapseMotion).Children(content)
		}
	})
}
