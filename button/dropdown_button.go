package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DropdownButton flips its chevron and opens a menu panel.
func DropdownButton(c *ui.Context, label string, items []string, open *bool, selected *int) ui.Element {
	if label == "" || len(items) == 0 || open == nil || selected == nil || *selected < -1 || *selected >= len(items) {
		panic("button: DropdownButton needs a label, items and state")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface).Expanded(*open)
	if root.Clicked() {
		*open = !*open
		root.SetChanged()
		if *open {
			sound.Play(sound.Open)
		} else {
			sound.Play(sound.Close)
		}
		kit.Log("DropdownButton", "toggled", "open", *open)
	}
	root.Children(func() {
		ui.Text(c, label).SingleLine()
		chevron := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
		chevron.Rotate(chevron.AnimateTo("flip", 180*unit(*open), ui.Snappy))
	})
	panel := ui.Popover(c, root, open, func() {
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			for i, item := range items {
				c.Key(i)
				row := ui.Button(c, item).FillWidth().Justify(ui.Start)
				if kit.Clicked(row, sound.Tap) {
					*selected = i
					*open = false
					root.SetChanged()
					kit.Log("DropdownButton", "selected", "index", i)
				}
			}
		})
	})
	if panel.Dismissed() {
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("DropdownButton", "dismissed")
	}
	return root
}
