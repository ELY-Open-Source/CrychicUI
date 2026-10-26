package choice

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Option is a radio whose dot scales in.
func Option(c *ui.Context, label, value string, selected *string) ui.Element {
	if selected == nil || label == "" {
		panic(fmt.Sprintf("choice: Option %q needs a label and selection", label))
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
