package choice

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChipChoice is one choice among pill chips.
func ChipChoice(c *ui.Context, label string, options []string, selected *string) ui.Element {
	distinct(who("ChipChoice", label), options)
	if label == "" || selected == nil || !slices.Contains(options, *selected) {
		panic("choice: " + who("ChipChoice", label) + " needs a label and a valid selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, option := range options {
			c.Key(option)
			chip := ui.RadioBase(c, selected, option).Label(option).Padding(t.Space(1), t.Space(3)).Radius(999)
			bg, fg := t.Surface, t.Text
			if option == *selected {
				bg, fg = t.Accent.Alpha(0.14), t.Accent
			} else if chip.Hovered() {
				bg = t.SurfaceHover
			}
			chip.Background(chip.AnimateColor("bg", bg, quick)).TextColor(chip.AnimateColor("fg", fg, quick))
			kit.Press(chip).Children(func() { ui.Text(c, option).SingleLine() })
			changed = kit.Changed(chip, sound.Select) || changed
		}
	}).Label(label).Row().Wrap().Gap(t.Space(2))
	if changed {
		root.SetChanged()
		kit.Log("ChipChoice", "selected", "index", slices.Index(options, *selected))
	}
	return root
}
