package files

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThumbsOptions name the zoom slider.
type ThumbsOptions struct{ Zoom string }

// Thumbs shows tiles sized by a zoom slider; cue Tick.
func Thumbs(c *ui.Context, label string, list []FileItem, zoom *float64, selected *int, o ThumbsOptions) ui.Element {
	validFiles("Thumbs", label, list)
	unique("Thumbs", label, list)
	if zoom == nil || selected == nil || !(*zoom >= 0.5 && *zoom <= 2) || math.IsInf(*zoom, 0) || *selected < -1 || *selected >= len(list) {
		panic(fmt.Sprintf("files: Thumbs %q needs zoom in [0.5,2] and a selection or -1", label))
	}
	def(&o.Zoom, "Zoom")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "images", label, func() {
			ui.Row(c).Width(t.Space(40)).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("zoom-out")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
				s := ui.Slider(c, zoom, 0.5, 2).Label(o.Zoom).Grow(1)
				if kit.Changed(s, sound.Tick) {
					root.SetChanged()
					kit.Log("Thumbs", "zoom", "value", *zoom)
				}
				ui.Icon(c, icon.Must("zoom-in")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
			})
		})
		side := root.AnimateTo("side", t.Space(24)*float32(*zoom), ui.Snappy)
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			for i, f := range list {
				c.Key(fileKey(f))
				on := i == *selected
				tile := kit.Press(ui.ButtonBase(c).Label(f.Name).Column().Gap(t.Space(1.5)).AlignItems(ui.Center).Width(side).
					Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, tile, i)
				tile.Background(tile.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(tile.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
				tile.Border(1, tile.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				tile.Children(func() {
					fileTile(c, f, side*0.55)
					ui.Text(c, f.Name).FontSize(t.Rem(0.75)).Bold().TextAlign(ui.Center).FillWidth().SingleLine()
				})
				if tile.Clicked() && !on {
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
