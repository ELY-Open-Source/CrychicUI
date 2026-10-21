package typography

import "github.com/egoist/mygo/ui"

// Title is bold text whose changes cross-fade.
func Title(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).Bold().FontSize(c.Theme().Rem(1.25)).ContentTransition(ui.ContentTransition{})
}
