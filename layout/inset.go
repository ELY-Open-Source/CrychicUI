package layout

import "github.com/egoist/mygo/ui"

// Inset is a recessed surface whose color follows the theme.
func Inset(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	e := ui.Column(c).Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius)
	return e.Background(e.AnimateColor("bg", t.Surface, ui.Smooth)).Children(content)
}
