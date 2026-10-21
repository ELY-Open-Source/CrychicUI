package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DropdownGroup is a titled section of options.
type DropdownGroup struct {
	Title   string
	Options []string
}

// GroupedDropdown is a select whose options sit under section titles.
func GroupedDropdown(c *ui.Context, label string, groups []DropdownGroup, selected *string) ui.Element {
	var all, titles []string
	for _, g := range groups {
		if len(g.Options) == 0 {
			panic("choice: GroupedDropdown groups need options")
		}
		titles, all = append(titles, g.Title), append(all, g.Options...)
	}
	distinct("GroupedDropdown", titles)
	distinct("GroupedDropdown", all)
	if label == "" || selected == nil || !slices.Contains(all, *selected) {
		panic("choice: GroupedDropdown needs a label and a valid selection")
	}
	t := c.Theme()
	sel := ui.SelectBase(c, selected)
	shown := *selected
	root := trigger(c, sel.Trigger.Label(label), sel.Open(), func() { ui.Text(c, shown).SingleLine() })
	sel.Popup(func(p ui.Element) {
		panel(c, p)
		for i, g := range groups {
			c.Key(g.Title)
			head := kit.Muted(c, g.Title).FontSize(t.Rem(0.8)).Bold().Padding(t.Space(1.5), t.Space(2.5), t.Space(1), t.Space(2.5))
			if i > 0 {
				head.BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Margin(t.Space(1), 0, 0, 0)
			}
			for _, option := range g.Options {
				item := sel.Item(option)
				entry(c, item, option, option == *selected, item.Highlighted())
			}
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("GroupedDropdown", "selected")
	}
	return root
}
