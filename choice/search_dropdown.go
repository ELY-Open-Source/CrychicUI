package choice

import (
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SearchDropdownOptions tune a SearchDropdown.
type SearchDropdownOptions struct {
	// Placeholder fills the empty filter; "Search" when empty.
	Placeholder string
	// Empty shows when nothing matches; "No results" when empty.
	Empty string
}

// SearchDropdown is a select whose popup filters as you type.
func SearchDropdown(c *ui.Context, label string, options []string, selected *string, o SearchDropdownOptions) ui.Element {
	distinct("SearchDropdown", options)
	if label == "" || selected == nil || !slices.Contains(options, *selected) {
		panic("choice: SearchDropdown needs a label and a valid selection")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search"
	}
	if o.Empty == "" {
		o.Empty = "No results"
	}
	t := c.Theme()
	root := ui.ButtonBase(c).Label(label)
	open := ui.Local(root, "open", func() bool { return false })
	query := ui.Local(root, "query", func() string { return "" })
	if root.Clicked() {
		*open, *query = !*open, ""
	}
	shown := *selected
	trigger(c, root, *open, func() { ui.Text(c, shown).SingleLine() })
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	ui.PopoverBase(c, root, open, func(p ui.Element) {
		panel(c, p).MinWidth(w).Gap(t.Space(1))
		var matches []string
		q := strings.ToLower(strings.TrimSpace(*query))
		for _, option := range options {
			if strings.Contains(strings.ToLower(option), q) {
				matches = append(matches, option)
			}
		}
		pick := ""
		ui.Row(c).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Icon(c, icon.Must("search")).Size(16, 16).TextColor(t.TextMuted)
			c.Key("filter")
			in := ui.TextInputBase(c, query).Grow(1).Label(o.Placeholder).Placeholder(o.Placeholder).AutoFocus()
			if in.Submitted() && len(matches) > 0 {
				pick = matches[0]
			}
		})
		if len(matches) == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(1.5), t.Space(2.5)).Transition(kit.Enter)
		}
		for _, option := range matches {
			c.Key(option)
			item := ui.ButtonBase(c).Label(option).Justify(ui.Start)
			entry(c, item, option, option == *selected, item.Hovered()).Transition(kit.Enter)
			if item.Clicked() {
				pick = option
			}
		}
		if pick != "" {
			*open = false
			sound.Play(sound.Select)
			if pick != *selected {
				*selected = pick
				root.SetChanged()
				kit.Log("SearchDropdown", "selected", "index", slices.Index(options, pick))
			}
		}
	})
	return root
}
