package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Check is a check box whose mark draws in.
func Check(c *ui.Context, label string, checked *bool) ui.Element {
	if checked == nil || label == "" {
		panic("choice: Check needs a label and state")
	}
	t := c.Theme()
	root := ui.CheckboxBase(c, checked).Label(label).Gap(t.Space(2))
	on := *checked
	root.Children(func() {
		box(c, root, on, false)
		ui.Text(c, label)
	})
	if kit.Changed(root, sound.Toggle) {
		kit.Log("Check", "toggled", "on", *checked)
	}
	return root
}
