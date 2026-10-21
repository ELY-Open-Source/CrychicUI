package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UploadButton accepts a dropped path and shows progress.
func UploadButton(c *ui.Context, label string, path *string, progress float64) ui.Element {
	if path == nil {
		panic("button: UploadButton needs a path")
	}
	root := progressButton(c, label, progress)
	if files := root.DroppedFiles(); len(files) > 0 {
		*path = files[0]
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("UploadButton", "dropped", "count", len(files))
	}
	return root
}
