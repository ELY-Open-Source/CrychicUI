package typography

import "github.com/egoist/mygo/ui"

// Truncate is one ellipsized line with a tooltip.
func Truncate(c *ui.Context, text string) ui.Element {
	return ui.Text(c, text).SingleLine().Tooltip(text).ContentTransition(ui.ContentTransition{})
}
