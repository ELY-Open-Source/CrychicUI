package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Option is a radio whose dot scales in.
func Option(c *ui.Context, label, value string, selected *string) ui.Element {
	if selected == nil || label == "" {
		panic("choice: Option needs a label and selection")
	}
	t := c.Theme()
	root := ui.RadioBase(c, selected, value).Label(label).Gap(t.Space(2))
	on := *selected == value
	root.Children(func() {
		dot(c, root, on)
		ui.Text(c, label)
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("Option", "selected")
	}
	return root
}
