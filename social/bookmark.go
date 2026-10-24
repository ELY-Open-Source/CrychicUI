package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
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
	label := o.Save
	if *saved {
		label = o.Saved
	}
	root := button(c, label).Row().Gap(c.Theme().Space(1)).AlignItems(ui.Center).Radius(999).Checked(*saved)
	root.Children(func() {
		ui.Icon(c, icon.Must("bookmark")).Size(16, 16)
		ui.Text(c, label).ContentTransition(ui.ContentTransition{})
	})
	if root.Clicked() {
		*saved = !*saved
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Bookmark", "saved", *saved)
	}
	return root
}
