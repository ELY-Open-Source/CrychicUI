package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Chip is an avatar and text pill lifting on hover.
func Chip(c *ui.Context, label string, avatar *ui.Bitmap) ui.Element {
	if avatar == nil {
		panic("display: Chip needs an avatar")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).Padding(t.Space(1), t.Space(2)).Radius(999).
		Background(t.Surface).Border(1, t.Border)
	kit.Lift(root).Children(func() {
		avatarCircle(c, label, avatar, t.Space(6))
		ui.Text(c, label)
	})
	return root
}
