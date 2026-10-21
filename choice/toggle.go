package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Toggle is a switch whose knob springs.
func Toggle(c *ui.Context, label string, on *bool) ui.Element {
	if on == nil || label == "" {
		panic("choice: Toggle needs a label and state")
	}
	t := c.Theme()
	root := ui.SwitchBase(c, on).Label(label).Gap(t.Space(2))
	shown := *on
	root.Children(func() {
		track(c, shown)
		ui.Text(c, label)
	})
	if kit.Changed(root, sound.Toggle) {
		kit.Log("Toggle", "toggled", "on", *on)
	}
	return root
}
