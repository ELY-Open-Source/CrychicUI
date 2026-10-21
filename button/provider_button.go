package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ProviderButton is a sign-in button with a provider icon.
func ProviderButton(c *ui.Context, label, iconName string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Justify(ui.Center)
	bg := t.Surface
	if root.Pressed() {
		bg = t.SurfacePressed
	} else if root.Hovered() {
		bg = t.SurfaceHover
	}
	root.Background(root.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(iconName)).Size(16, 16)
		ui.Text(c, label).SingleLine()
	})
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root
}
