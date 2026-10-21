package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// InfoCard is a card with title, body and footer slots.
func InfoCard(c *ui.Context, title, body string, footer func()) ui.Element {
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4))
	kit.Lift(root).Children(func() {
		ui.Text(c, title).Bold()
		ui.Text(c, body).TextColor(t.TextMuted)
		if footer != nil {
			ui.Row(c).Padding(t.Space(2), 0, 0, 0).Children(footer)
		}
	})
	return root
}
