package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LightboxOptions tune a Lightbox.
type LightboxOptions struct {
	Back, Forth, Close string
}

// Lightbox shows images over a dark backdrop with arrows.
func Lightbox(c *ui.Context, label string, images []*ui.Bitmap, captions []string, at *int, open *bool, o LightboxOptions) ui.Element {
	if label == "" || len(images) == 0 || at == nil || open == nil || !valid(*at, len(images)) || (captions != nil && len(captions) != len(images)) {
		panic(fmt.Sprintf("overlay: Lightbox %q needs images, a caption each or none, a shown one and open", label))
	}
	for _, img := range images {
		if img == nil {
			panic(fmt.Sprintf("overlay: Lightbox %q images must not be nil", label))
		}
	}
	def(&o.Back, "Previous")
	def(&o.Forth, "Next")
	def(&o.Close, "Close")
	t := c.Theme()
	root := ui.Box(c)
	shown, last := *at, len(images)-1
	back, forth := "chevron-left", "chevron-right"
	prev, next := ui.KeyLeft, ui.KeyRight
	if c.Direction() == ui.RightToLeft {
		back, forth, prev, next = forth, back, next, prev
	}
	ink := lighter(t.Text, t.Background)
	turn := func(to int) {
		*at = to
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("Lightbox", "shown", "index", to)
	}
	shut := func() {
		*open = false
		root.SetChanged()
	}
	arrow := func(name, glyph string, off bool) bool {
		b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(999)).Disabled(off)
		b.Background(b.AnimateColor("bg", ink.Alpha(0.12+0.12*unit(b.Hovered())), ui.Snappy)).TextColor(ink).
			Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
		return b.Clicked()
	}
	if *open {
		ui.Overlay(c, func() {
			scrim := ui.Column(c).Absolute().Left(0).Top(0).Right(0).Bottom(0).Justify(ui.Center).AlignItems(ui.Center).
				Gap(t.Space(3)).Padding(t.Space(4)).Modal().Role(ui.RoleDialog).Label(label)
			kit.Arrive(scrim.Background(kit.Shade(0.88)), ui.Pose{Fade: 1}, 0).Children(func() {
				ui.Row(c).FillWidth().Justify(ui.End).Children(func() {
					if arrow(o.Close, "x", false) {
						shut()
					}
				})
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
					if arrow(o.Back, back, shown == 0) {
						turn(shown - 1)
					}
					ui.Box(c).Grow(1).MaxWidth(t.Space(140)).Children(func() {
						kit.Slide(c, shown, func() {
							ui.Image(c, images[shown]).FillWidth().AspectRatio(1.4)
						})
					})
					if arrow(o.Forth, forth, shown == last) {
						turn(shown + 1)
					}
				})
				if captions != nil {
					ui.Text(c, captions[shown]).TextColor(ink).ContentTransition(ui.ContentTransition{})
				}
			})
			switch {
			case scrim.OverlayShortcut(0, ui.KeyEscape):
				shut()
			case scrim.OverlayShortcut(0, prev) && shown > 0:
				turn(shown - 1)
			case scrim.OverlayShortcut(0, next) && shown < last:
				turn(shown + 1)
			}
		})
	}
	moved(root, *open, "Lightbox")
	return root
}
