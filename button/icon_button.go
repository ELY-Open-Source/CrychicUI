package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IconButton is a square icon button with a required label.
func IconButton(c *ui.Context, label, iconName string) ui.Element {
	if label == "" {
		panic("button: IconButton needs a label")
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius))
	bg := t.Surface
	if b.Pressed() {
		bg = t.SurfacePressed
	} else if b.Hovered() {
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Border(1, t.Border).Children(func() {
		ui.Icon(c, icon.Must(iconName)).Size(16, 16)
	})
	kit.Clicked(b, sound.Tap)
	return b
}
