package media

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MediaItem is a photo, or a clip with Length.
type MediaItem struct {
	Name   string
	Image  *ui.Bitmap
	Length time.Duration
}

// MediaGridOptions tune a MediaGrid.
type MediaGridOptions struct {
	// Columns per row; zero gives 3.
	Columns int
	// Photo and Clip name the badges.
	Photo, Clip string
}

// MediaGrid shows photos and clips with type badges.
func MediaGrid(c *ui.Context, label string, items []MediaItem, picked *int, o MediaGridOptions) ui.Element {
	if label == "" || len(items) == 0 || o.Columns < 0 {
		panic(fmt.Sprintf("media: MediaGrid %q needs a label, items and columns not negative", label))
	}
	for i, it := range items {
		if it.Name == "" || it.Image == nil || it.Length < 0 || slices.IndexFunc(items, func(m MediaItem) bool { return m.Name == it.Name }) != i {
			panic(fmt.Sprintf("media: MediaGrid %q items need distinct names, images and lengths not negative", label))
		}
	}
	chosen("MediaGrid", label, picked, -1, len(items))
	if o.Columns == 0 {
		o.Columns = 3
	}
	def(&o.Photo, "Photo")
	def(&o.Clip, "Video")
	t := c.Theme()
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	shown := *picked
	return root.Children(func() {
		for r := 0; r < len(items); r += o.Columns {
			c.Key(r)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for i := r; i < r+o.Columns; i++ {
					c.Key(i)
					if i >= len(items) {
						ui.Box(c).Grow(1).Basis(0)
						continue
					}
					it := items[i]
					kind, glyph, note := o.Photo, "image", o.Photo
					if it.Length > 0 {
						kind, glyph, note = o.Clip, "video", clock(it.Length)
					}
					ring := t.Accent.Alpha(unit(i == shown))
					tile := kit.Press(ui.ButtonBase(c).Label(it.Name + ", " + kind).Tooltip(it.Name).Grow(1).Basis(0).AspectRatio(1).
						Radius(t.Radius).Clip().Background(t.Surface))
					tile.Checked(i == shown).Border(t.Space(0.5), tile.AnimateColor("ring", ring, ui.Snappy))
					kit.Arrive(tile, ui.Pose{Fade: 1, Y: t.Space(2)}, time.Duration(i)*40*time.Millisecond)
					tile.Children(func() {
						ui.Grid(c).Columns(1).Fill().Children(func() {
							pose := ui.Pose{}
							if tile.Hovered() {
								pose.Scale = 1.06
							}
							ui.Image(c, it.Image).ColumnStart(1).RowStart(1).Fill().Fit(ui.Cover).AnimatePose(pose, ui.Snappy)
							ui.Row(c).ColumnStart(1).RowStart(1).AlignSelf(ui.Start).JustifySelf(ui.Start).Margin(t.Space(2)).
								Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
								Background(kit.Shade(0.6)).TextColor(ink).Children(func() {
								ui.Icon(c, icon.Must(glyph)).Size(16, 16)
								kit.Numeric(c, note).FontSize(t.Rem(0.75))
							})
						})
					})
					if tile.Clicked() && i != *picked {
						*picked = i
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("MediaGrid", "picked", "index", i)
					}
				}
			})
		}
	})
}
