package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Dropdown is a select whose popup scales in.
func Dropdown(c *ui.Context, label string, options []string, selected *string) ui.Element {
	distinct(who("Dropdown", label), options)
	if label == "" || selected == nil || !slices.Contains(options, *selected) {
		panic("choice: " + who("Dropdown", label) + " needs a label and a valid selection")
	}
	sel := ui.SelectBase(c, selected)
	shown := *selected
	root := trigger(c, sel.Trigger.Label(label), sel.Open(), func() { ui.Text(c, shown).SingleLine() })
	sel.Popup(func(p ui.Element) {
		panel(c, p)
		optionList(c, func() {
			for _, option := range options {
				c.Key(option)
				item := sel.Item(option)
				entry(c, item, option, option == *selected, item.Highlighted())
			}
		})
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("Dropdown", "selected", "index", slices.Index(options, *selected))
	}
	return root
}
