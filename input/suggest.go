package input

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Suggest is a text input with filtering styled suggestions.
func Suggest(c *ui.Context, label string, value *string, suggestions []string) ui.Element {
	if value == nil || len(suggestions) == 0 {
		panic("input: Suggest needs a value and suggestions")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	open := ui.Local(root, "open", func() bool { return false })
	dismissed := ui.Local(root, "dismissed", func() string { return "" })
	lastQuery := ui.Local(root, "lastQuery", func() string { return *value })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInput(c, value).FillWidth().Placeholder(label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		changed(root, in)
		var matches []string
		for _, item := range suggestions {
			if strings.Contains(strings.ToLower(item), strings.ToLower(*value)) {
				matches = append(matches, item)
			}
		}
		if *value != *lastQuery {
			*dismissed = ""
			*lastQuery = *value
		}
		if in.Focused() && len(matches) > 0 && *dismissed != *value {
			*open = true
		}
		panel := ui.Popover(c, in, open, func() {
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				for i, item := range matches {
					c.Key(i)
					row := ui.Button(c, item).FillWidth().Justify(ui.Start)
					if kit.Clicked(row, sound.Select) {
						*value = item
						*lastQuery = item
						*open = false
						*dismissed = item
						root.SetChanged()
						kit.Log("Suggest", "selected", "index", i)
					}
				}
			})
		})
		if panel.Dismissed() {
			*open = false
			*dismissed = *value
		}
	})
}
