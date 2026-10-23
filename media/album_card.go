package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AlbumCardOptions tune an AlbumCard.
type AlbumCardOptions struct {
	Play string
	// Width of the card; zero gives the default.
	Width float32
}

// AlbumCard is a cover whose play button rises.
func AlbumCard(c *ui.Context, title, artist string, art *ui.Bitmap, playing *bool, o AlbumCardOptions) ui.Element {
	if title == "" || playing == nil {
		panic(fmt.Sprintf("media: AlbumCard %q needs a title and state", title))
	}
	def(&o.Play, "Play")
	t := c.Theme()
	size("AlbumCard", title, &o.Width, t.Space(44))
	root := kit.Lift(ui.Column(c).Width(o.Width).Gap(t.Space(2)).Label(title))
	return root.Children(func() {
		ui.Grid(c).Columns(1).FillWidth().AspectRatio(1).Radius(t.Radius).Clip().Background(t.Surface).Children(func() {
			cover(c, art, o.Width).ColumnStart(1).RowStart(1)
			wrap := ui.Box(c).ColumnStart(1).RowStart(1).AlignSelf(ui.End).JustifySelf(ui.End).Padding(t.Space(2))
			wrap.Children(func() {
				play := button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play})
				if play.Changed() {
					root.SetChanged()
					kit.Log("AlbumCard", "toggled", "playing", *playing)
				}
				up := unit(root.Hovered() || play.Focused() || *playing)
				wrap.AnimatePose(ui.Pose{Y: t.Space(3) * (1 - up), Fade: 1 - up}, ui.Snappy)
			})
		})
		ui.Text(c, title).SingleLine().Bold()
		kit.Muted(c, artist).SingleLine().FontSize(t.Rem(0.8))
	})
}
