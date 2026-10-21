package choice

import (
	"slices"

	"github.com/egoist/mygo/ui"
)

// OptionGroupOptions tune an OptionGroup.
type OptionGroupOptions struct {
	// Row lays the options out side by side.
	Row bool
}

// OptionGroup is radios in a row or column.
func OptionGroup(c *ui.Context, label string, options []string, selected *string, o OptionGroupOptions) ui.Element {
	distinct("OptionGroup", options)
	if label == "" || selected == nil || !slices.Contains(options, *selected) {
		panic("choice: OptionGroup needs a label and a valid selection")
	}
	t := c.Theme()
	changed := false
	root := ui.RadioGroup(c, func() {
		for _, option := range options {
			c.Key(option)
			changed = Option(c, option, option, selected).Changed() || changed
		}
	}).Label(label)
	if changed {
		root.SetChanged()
	}
	if o.Row {
		root.Row().Gap(t.Space(4))
	}
	return root
}
