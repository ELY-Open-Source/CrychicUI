package navigation

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// tabs checks a tab component's labels and active index.
func tabs(component, label string, names []string, active *int) {
	if label == "" || len(names) == 0 || active == nil || !valid(*active, len(names)) {
		panic(fmt.Sprintf("navigation: %s needs a label, tabs and an active index", component))
	}
}

// underlined styles one tab; extra, if set, follows the name.
func underlined(c *ui.Context, tab ui.Element, name string, on bool, id uint64, extra func()) ui.Element {
	t := c.Theme()
	fg := t.TextMuted
	if on || tab.Hovered() {
		fg = t.Text
	}
	return tab.Padding(t.Space(2), t.Space(3), 0, t.Space(3)).TextColor(tab.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
		ui.Column(c).Gap(t.Space(1.5)).AlignItems(ui.Stretch).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, name).SingleLine()
				if extra != nil {
					extra()
				}
			})
			line := ui.Box(c).Height(2).Radius(999)
			if on {
				line.Background(t.Accent).Transition(glide(id))
			}
		})
	})
}
