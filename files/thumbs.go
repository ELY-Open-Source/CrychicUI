package files

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThumbsOptions are the thumbnail labels.
type ThumbsOptions struct{ Zoom string }

// Thumbs shows large icons with a zoom slider; cue Tick.
func Thumbs(c *ui.Context, label string, list []FileItem, zoom *float64, selected *int, o ThumbsOptions) ui.Element {
	validFiles("Thumbs", label, list)
	if zoom == nil || selected == nil || *zoom < 0.5 || *zoom > 2 || math.IsNaN(*zoom) || math.IsInf(*zoom, 0) ||
		*selected < -1 || *selected >= len(list) {
		panic("files: Thumbs needs valid zoom and selection")
	}
	t := c.Theme()
	def(&o.Zoom, "Zoom")
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		s := ui.Slider(c, zoom, 0.5, 2).Label(o.Zoom).FillWidth()
		if kit.Changed(s, sound.Tick) {
			root.SetChanged()
			kit.Log("Thumbs", "zoom", "value", *zoom)
		}
		ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
			for i, item := range list {
				c.Key(i)
				side := t.Space(12) * float32(*zoom)
				tile := kit.Press(ui.ButtonBase(c).Label(item.Name).Column().Gap(t.Space(1)).AlignItems(ui.Center).
					Width(side).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Checked(i == *selected))
				tile.Background(tile.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
				tile.Children(func() {
					ui.Icon(c, icon.Must(itemGlyph(item))).Size(side/2, side/2).TextColor(t.Accent)
					ui.Text(c, item.Name).FontSize(t.Rem(0.75)).SingleLine()
				})
				if tile.Clicked() && *selected != i {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Thumbs", "selected", "index", i)
				}
			}
		})
	})
	return root
}
