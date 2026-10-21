package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToggleRow is a setting row with title, description and switch.
func ToggleRow(c *ui.Context, title, description string, on *bool) ui.Element {
	if on == nil || title == "" {
		panic("choice: ToggleRow needs a title and state")
	}
	t := c.Theme()
	root := ui.SwitchBase(c, on).Label(title).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), 0)
	shown := *on
	root.Children(func() {
		ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, title)
			if description != "" {
				kit.Muted(c, description)
			}
		})
		track(c, shown)
	})
	if kit.Changed(root, sound.Toggle) {
		kit.Log("ToggleRow", "toggled", "on", *on)
	}
	return root
}
