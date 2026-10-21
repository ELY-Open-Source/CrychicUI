package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// IconBoxOptions tune an IconBox.
type IconBoxOptions struct {
	// Color tints the box and icon; zero is Accent.
	Color ui.Color
}

// IconBox is an icon in a tinted rounded square.
func IconBox(c *ui.Context, name string, o IconBoxOptions) ui.Element {
	if o.Color.A == 0 {
		o.Color = c.Theme().Accent
	}
	t := c.Theme()
	box := ui.Column(c).AlignSelf(ui.Start).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius)
	box.Background(box.AnimateColor("tint", o.Color.Alpha(0.14), ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(name)).Size(16, 16).TextColor(box.AnimateColor("icon", o.Color, ui.Snappy))
	})
	return box
}
