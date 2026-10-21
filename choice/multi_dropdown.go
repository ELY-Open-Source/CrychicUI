package choice

import (
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MultiDropdown is a select of checks with a count.
func MultiDropdown(c *ui.Context, label string, options []string, selected *[]string) ui.Element {
	distinct("MultiDropdown", options)
	subset("MultiDropdown", options, selected)
	if label == "" {
		panic("choice: MultiDropdown needs a label")
	}
	t := c.Theme()
	root := ui.ButtonBase(c).Label(label)
	open := ui.Local(root, "open", func() bool { return false })
	if root.Clicked() {
		*open = !*open
	}
	n := len(*selected)
	trigger(c, root, *open, func() {
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			ui.Text(c, label).SingleLine()
			count := kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.85)).Padding(0, t.Space(1.5)).Radius(999).
				Background(t.Accent).TextColor(t.AccentText)
			k := count.AnimateTo("shown", unit(n > 0), ui.Bouncy)
			count.Scale(k, k).Opacity(min(max(k, 0), 1))
		})
	})
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	ui.PopoverBase(c, root, open, func(p ui.Element) {
		panel(c, p).MinWidth(w)
		for _, option := range options {
			on := slices.Contains(*selected, option)
			c.Key(option)
			row := ui.CheckboxBase(c, &on).Label(option).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)
			row.Background(row.AnimateColor("bg", t.Surface.Alpha(unit(row.Hovered())), quick))
			shown := on
			row.Children(func() {
				box(c, row, shown, false)
				ui.Text(c, option).SingleLine()
			})
			if kit.Changed(row, sound.Toggle) {
				*selected = flip(*selected, option)
				root.SetChanged()
				kit.Log("MultiDropdown", "toggled", "count", len(*selected))
			}
		}
	})
	return root
}
