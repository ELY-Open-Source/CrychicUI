package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MediaCard is a card whose top image zooms on hover.
func MediaCard(c *ui.Context, title, body string, image *ui.Bitmap) ui.Element {
	if image == nil {
		panic("display: MediaCard needs an image")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Clip()
	return root.Children(func() {
		img := ui.Image(c, image).FillWidth().Height(t.Rem(8)).Fit(ui.Cover)
		scale := img.AnimateTo("zoom", 1+0.05*unit(root.Hovered()), ui.Snappy)
		img.Scale(scale, scale)
		ui.Column(c).Gap(t.Space(1)).Padding(t.Space(3)).Children(func() {
			ui.Text(c, title).Bold()
			kit.Muted(c, body)
		})
	})
}
