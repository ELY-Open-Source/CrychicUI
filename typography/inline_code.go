package typography

import "github.com/egoist/mygo/ui"

// InlineCode is monospace text on a surface pill.
func InlineCode(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	return ui.Text(c, text).AlignSelf(ui.Start).Font("monospace").FontSize(t.Rem(0.85)).Padding(t.Space(0.5), t.Space(1)).
		Radius(999).Background(t.Surface).ContentTransition(ui.ContentTransition{})
}
