package typography

import "github.com/egoist/mygo/ui"

// Paragraph is wrapping body text with cross-fading changes.
func Paragraph(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).ContentTransition(ui.ContentTransition{})
}
