package layout

import "github.com/egoist/mygo/ui"

// Panel is a bordered surface box; its size animates.
func Panel(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).Gap(t.Space(2)).Padding(t.Space(4)).Background(t.Surface).Border(1, t.Border).
		Radius(t.Radius).Transition(ui.ElementTransition{Size: true}).Children(content)
}
