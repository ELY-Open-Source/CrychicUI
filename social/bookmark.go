package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BookmarkOptions are the bookmark labels.
type BookmarkOptions struct{ Save, Saved string }

// Bookmark toggles saved state; cue Toggle.
func Bookmark(c *ui.Context, saved *bool, o BookmarkOptions) ui.Element {
	if saved == nil {
		panic("social: Bookmark needs state")
	}
	def(&o.Save, "Save")
	def(&o.Saved, "Saved")
	t := c.Theme()
	label, mark := o.Save, "bookmark"
	if *saved {
		label, mark = o.Saved, "bookmark-check"
	}
	root := kit.Press(ui.ButtonBase(c).Label(label).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).AlignSelf(ui.Start).
		Padding(t.Space(1.5), t.Space(3)).Radius(999).Checked(*saved))
	hover := unit(root.Hovered())
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, hover).Mix(t.Accent.Alpha(0.12), unit(*saved)), ui.Snappy))
	root.Border(1, root.AnimateColor("line", t.Border.Mix(t.Accent.Alpha(0.5), unit(*saved)), ui.Snappy))
	root.TextColor(root.AnimateColor("ink", t.Text.Mix(t.Accent, unit(*saved)), ui.Snappy))
	root.Children(func() {
		c.Key(*saved)
		icon := glyph(c, mark, 16)
		if *saved {
			kit.Arrive(icon, ui.Pose{Y: -t.Space(2), Scale: 1.4, Fade: 1}, 0)
		}
		ui.Text(c, label).FontSize(t.Rem(0.82)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*saved = !*saved
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Bookmark", "saved", *saved)
	}
	return root
}
