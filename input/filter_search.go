package input

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FilterSearchOptions tune a FilterSearch.
type FilterSearchOptions struct {
	// Placeholder is caller copy.
	Placeholder string
}

// FilterSearch has a search field with toggleable filter chips.
func FilterSearch(c *ui.Context, query *string, filters []string, selected *[]string, o FilterSearchOptions) ui.Element {
	if query == nil || selected == nil || len(filters) == 0 {
		panic("input: FilterSearch needs query, filters and selection")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		in := ui.TextInput(c, query).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		changed(root, in)
		ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
			for i, filter := range filters {
				on := slices.Contains(*selected, filter)
				chip := ui.ButtonBase(c).Key(i).Label(filter).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(999).
					Checked(on).Border(1, t.Border).Transition(ui.ElementTransition{Position: true, Size: true})
				bg, fg := t.Surface, t.Text
				if on {
					bg, fg = t.Accent.Alpha(0.14), t.Accent
				}
				chip.Background(chip.AnimateColor("bg", bg, ui.Snappy)).TextColor(fg).Children(func() {
					if on {
						ui.Icon(c, icon.Must("check")).Size(16, 16).ContentTransition(ui.ContentTransition{})
					}
					ui.Text(c, filter)
				})
				if kit.Clicked(chip, sound.Toggle) {
					if on {
						*selected = slices.DeleteFunc(*selected, func(s string) bool { return s == filter })
					} else {
						*selected = append(*selected, filter)
					}
					root.SetChanged()
					kit.Log("FilterSearch", "toggled", "count", len(*selected))
				}
			}
		})
	})
}
