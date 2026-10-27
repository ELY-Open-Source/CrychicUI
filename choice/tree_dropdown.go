package choice

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TreeItem is a labeled item with nested children.
type TreeItem struct {
	Label    string
	Children []TreeItem
}

// treeLabels lists every label of items, depth first.
func treeLabels(items []TreeItem) []string {
	var out []string
	for _, it := range items {
		out = append(out, it.Label)
		out = append(out, treeLabels(it.Children)...)
	}
	return out
}

// TreeDropdownOptions tune a TreeDropdown; defaults read "Expand %s".
type TreeDropdownOptions struct {
	// Expand and Collapse format the branch button name.
	Expand, Collapse string
}

// TreeDropdown is a select of nested items that expand.
func TreeDropdown(c *ui.Context, label string, items []TreeItem, selected *string, o TreeDropdownOptions) ui.Element {
	labels := treeLabels(items)
	distinct(who("TreeDropdown", label), labels)
	if label == "" || selected == nil || !slices.Contains(labels, *selected) {
		panic("choice: " + who("TreeDropdown", label) + " needs a label and a valid selection")
	}
	if o.Expand == "" {
		o.Expand = "Expand %s"
	}
	if o.Collapse == "" {
		o.Collapse = "Collapse %s"
	}
	t := c.Theme()
	sel := ui.SelectBase(c, selected)
	shown := *selected
	root := trigger(c, sel.Trigger.Label(label), sel.Open(), func() { ui.Text(c, shown).SingleLine() })
	expanded := ui.Local(root, "expanded", func() map[string]bool { return map[string]bool{} })
	var level func(items []TreeItem, depth int)
	level = func(items []TreeItem, depth int) {
		for _, it := range items {
			c.Key(it.Label)
			ui.Row(c).PaddingStart(t.Space(5) * float32(depth)).Transition(kit.Enter).Children(func() {
				if len(it.Children) == 0 {
					ui.Box(c).Size(t.Space(6), t.Space(6))
				} else {
					name := o.Expand
					if (*expanded)[it.Label] {
						name = o.Collapse
					}
					expand := ui.ButtonBase(c).Label(fmt.Sprintf(name, it.Label)).Expanded((*expanded)[it.Label]).
						Size(t.Space(6), t.Space(6)).Radius(t.Radius).TextColor(t.TextMuted)
					kit.Press(expand).Children(func() {
						chevron := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16)
						chevron.Rotate(chevron.AnimateTo("turn", 90*unit((*expanded)[it.Label]), ui.Snappy))
					})
					if kit.Clicked(expand, sound.Tick) {
						(*expanded)[it.Label] = !(*expanded)[it.Label]
						kit.Log("TreeDropdown", "expanded", "open", (*expanded)[it.Label])
					}
				}
				item := sel.Item(it.Label).Grow(1)
				entry(c, item, it.Label, it.Label == *selected, item.Highlighted())
			})
			if (*expanded)[it.Label] {
				level(it.Children, depth+1)
			}
		}
	}
	sel.Popup(func(p ui.Element) {
		panel(c, p)
		optionList(c, func() { level(items, 0) })
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("TreeDropdown", "selected")
	}
	return root
}
