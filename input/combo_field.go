package input

import (
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ComboField is a styled combobox with a scaling popup.
func ComboField(c *ui.Context, label string, selected *string, options []string) ui.Element {
	if selected == nil || !slices.Contains(options, *selected) {
		panic("input: ComboField needs a selected option")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth()
	text := ui.Local(root, "text", func() string { return *selected })
	return root.Children(func() {
		c.Key("combo")
		combo := ui.ComboboxBase(c, text)
		combo.Input.FillWidth().Label(label).Placeholder(label).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border)
		if !combo.Input.Focused() && *text != *selected {
			*text = *selected
		}
		combo.Popup(func(panel ui.Element) {
			panel.Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
			kit.Arrive(panel, ui.Pose{Scale: 0.96, Fade: 1}, 0)
			shown := options
			if combo.Filtering() {
				shown = matchingOptions(options, *text)
			}
			for _, option := range shown {
				item := combo.Item(option).Padding(t.Space(2), t.Space(3)).Radius(t.Radius - t.Space(0.5))
				if item.Highlighted() {
					item.Background(t.Accent).TextColor(t.AccentText)
				}
				item.Children(func() { ui.Text(c, option) })
			}
		})
		if chosen, ok := combo.Chosen(); ok {
			*selected = chosen
			*text = chosen
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("ComboField", "selected", "index", slices.Index(options, chosen))
		}
	})
}

// matchingOptions filters options, prefix matches first.
func matchingOptions(options []string, query string) []string {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return options
	}
	var starts, contains []string
	for _, option := range options {
		lower := strings.ToLower(option)
		if strings.HasPrefix(lower, q) {
			starts = append(starts, option)
		} else if strings.Contains(lower, q) {
			contains = append(contains, option)
		}
	}
	return append(starts, contains...)
}
