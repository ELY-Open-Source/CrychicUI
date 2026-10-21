package typography

import (
	"strings"

	"github.com/egoist/mygo/ui"
)

// Overline is an uppercase spaced label with cross-fading changes.
func Overline(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	return ui.Text(c, strings.ToUpper(text)).TextColor(t.TextMuted).FontSize(t.Rem(0.75)).
		LetterSpacing(t.Rem(0.08)).ContentTransition(ui.ContentTransition{})
}
