package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// GalleryOptions tune a Gallery.
type GalleryOptions struct {
	// Columns per row; zero gives 3.
	Columns int
	// Item names a tile by number; "Photo %d" by default.
	Item               string
	Back, Forth, Close string
}

// Gallery is an image grid; open shows one large.
func Gallery(c *ui.Context, label string, images []*ui.Bitmap, open *int, o GalleryOptions) ui.Element {
	pictures("Gallery", label, images)
	chosen("Gallery", label, open, -1, len(images))
	if o.Columns < 0 {
		panic(fmt.Sprintf("media: Gallery %q columns must not be negative", label))
	}
	if o.Columns == 0 {
		o.Columns = 3
	}
	def(&o.Item, "Photo %d")
	def(&o.Back, "Previous")
	def(&o.Forth, "Next")
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Label(label)
	shown, last := *open, len(images)-1
	set := func(to int, cue sound.Cue, event string) {
		*open = to
		root.SetChanged()
		sound.Play(cue)
		kit.Log("Gallery", event, "index", to)
	}
	root.Children(func() {
		for r := 0; r <= last; r += o.Columns {
			c.Key(r)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for i := r; i < r+o.Columns; i++ {
					c.Key(i)
					if i > last {
						ui.Box(c).Grow(1).Basis(0)
						continue
					}
					name := fmt.Sprintf(o.Item, i+1)
					tile := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Grow(1).Basis(0).AspectRatio(1).
						Radius(t.Radius).Clip().Background(t.Surface))
					tile.Children(func() {
						pose := ui.Pose{}
						if tile.Hovered() {
							pose.Scale = 1.06
						}
						ui.Image(c, images[i]).Fill().Fit(ui.Cover).AnimatePose(pose, ui.Snappy)
					})
					kit.Arrive(tile, ui.Pose{Fade: 1, Y: t.Space(2)}, time.Duration(i)*40*time.Millisecond)
					if tile.Clicked() {
						set(i, sound.Open, "opened")
					}
				}
			})
		}
	})
	if shown < 0 {
		return root
	}
	back, forth, prev, next := ends(c)
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	arrow := func(name, glyph string, off bool) bool {
		b := round(c, name, glyph, t.Space(11)).Disabled(off)
		b.Background(b.AnimateColor("plate", ink.Alpha(0.12+0.12*unit(b.Hovered())), ui.Snappy)).TextColor(ink)
		return b.Clicked()
	}
	ui.Overlay(c, func() {
		scrim := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Justify(ui.Center).AlignItems(ui.Center).
			Gap(t.Space(3)).Padding(t.Space(4)).Modal().Role(ui.RoleDialog).Label(label)
		kit.Arrive(scrim.Background(kit.Shade(0.88)), ui.Pose{Fade: 1}, 0).Children(func() {
			ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
				if arrow(o.Close, "x", false) {
					set(-1, sound.Close, "closed")
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
				if arrow(o.Back, back, shown == 0) {
					set(shown-1, sound.Swipe, "shown")
				}
				ui.Box(c).Grow(1).MaxWidth(t.Space(140)).Children(func() {
					kit.Slide(c, shown, func() {
						ui.Image(c, images[shown]).FillWidth().AspectRatio(1.4).Label(fmt.Sprintf(o.Item, shown+1))
					})
				})
				if arrow(o.Forth, forth, shown == last) {
					set(shown+1, sound.Swipe, "shown")
				}
			})
		})
		switch {
		case scrim.OverlayShortcut(0, ui.KeyEscape):
			set(-1, sound.Close, "closed")
		case scrim.OverlayShortcut(0, prev) && shown > 0:
			set(shown-1, sound.Swipe, "shown")
		case scrim.OverlayShortcut(0, next) && shown < last:
			set(shown+1, sound.Swipe, "shown")
		}
	})
	return root
}
