package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ActionBarOptions tune an ActionBar.
type ActionBarOptions struct {
	// Save and Cancel are button copy.
	Save, Cancel string
}

// ActionBar scrolls content above a Save bar shown when dirty.
func ActionBar(c *ui.Context, dirty bool, o ActionBarOptions, content func()) (root, save, cancel ui.Element) {
	if content == nil {
		panic("form: ActionBar needs content")
	}
	def(&o.Save, "Save")
	def(&o.Cancel, "Cancel")
	root = ui.Column(c).Fill().Clip()
	root.Children(func() {
		layout.ScrollArea(c, content).Grow(1)
		cancel, save = bar(c, dirty, func() {}, o.Cancel, o.Save)
	})
	if kit.Clicked(save, sound.Success) {
		kit.Log("ActionBar", "save")
	}
	if kit.Clicked(cancel, sound.Tap) {
		kit.Log("ActionBar", "cancel")
	}
	return root, save, cancel
}
