package form

import (
	"path/filepath"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileFieldOptions tune a FileField.
type FileFieldOptions struct {
	// Empty, Browse and Remove are caller copy.
	Empty, Browse, Remove string
}

// FileField takes a dropped file; browse opens the caller's picker.
func FileField(c *ui.Context, label string, path *string, o FileFieldOptions) (root, browse ui.Element) {
	if label == "" || path == nil {
		panic("form: FileField needs a label and a path")
	}
	def(&o.Empty, "Drop a file here")
	def(&o.Browse, "Browse")
	def(&o.Remove, "Remove file")
	t := c.Theme()
	root = ui.Row(c).Label(label).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius)
	border := t.Border
	if root.FileDragOver() {
		border = t.Accent
	}
	root.Border(1, root.AnimateColor("border", border, ui.Snappy))
	take(root, path, root.DroppedFiles(), "FileField")
	root.Children(func() {
		ui.Icon(c, icon.Must("paperclip")).Size(16, 16).TextColor(t.TextMuted)
		if *path == "" {
			ui.Text(c, o.Empty).Grow(1).SingleLine().TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{})
		} else {
			ui.Text(c, filepath.Base(*path)).Grow(1).SingleLine().ContentTransition(ui.ContentTransition{})
		}
		if *path != "" && cross(c, o.Remove).Clicked() {
			*path = ""
			root.SetChanged()
			sound.Play(sound.Delete)
			kit.Log("FileField", "removed")
		}
		browse = plain(c, o.Browse)
	})
	kit.Clicked(browse, sound.Open)
	return root, browse
}
