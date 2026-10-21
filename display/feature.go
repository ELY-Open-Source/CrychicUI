package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Feature is an icon box, title and text rising in.
func Feature(c *ui.Context, iconName, title, text string) ui.Element {
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Children(func() {
		IconBox(c, iconName, IconBoxOptions{})
		ui.Text(c, title).Bold()
		kit.Muted(c, text)
	})
	return kit.Arrive(root, ui.Pose{Y: 8, Fade: 1}, 0)
}
