package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DirtyBarOptions tune a DirtyBar.
type DirtyBarOptions struct {
	// Message, Discard and Save are caller copy.
	Message, Discard, Save string
}

// DirtyBar slides up with Discard and Save while dirty.
func DirtyBar(c *ui.Context, dirty bool, o DirtyBarOptions) (root, discard, save ui.Element) {
	def(&o.Message, "Unsaved changes")
	def(&o.Discard, "Discard")
	def(&o.Save, "Save")
	root = ui.Column(c).FillWidth().Clip()
	root.Children(func() {
		discard, save = bar(c, dirty, func() { ui.Text(c, o.Message).FontWeight(500) }, o.Discard, o.Save)
	})
	if kit.Clicked(discard, sound.Close) {
		kit.Log("DirtyBar", "discard")
	}
	if kit.Clicked(save, sound.Success) {
		kit.Log("DirtyBar", "save")
	}
	return root, discard, save
}
