package form

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageFieldOptions tune an ImageField.
type ImageFieldOptions struct {
	// Empty, Choose, Replace and Remove are caller copy.
	Empty, Choose, Replace, Remove string
}

// ImageField previews the caller's decoded image; pick opens its picker.
func ImageField(c *ui.Context, label string, path *string, preview *ui.Bitmap, o ImageFieldOptions) (root, pick ui.Element) {
	if label == "" || path == nil {
		panic("form: ImageField needs a label and a path")
	}
	def(&o.Empty, "Drop an image")
	def(&o.Choose, "Choose image")
	def(&o.Replace, "Replace")
	def(&o.Remove, "Remove image")
	t := c.Theme()
	root = ui.Row(c).Label(label).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
	take(root, path, root.DroppedFiles(), "ImageField")
	root.Children(func() {
		frame := ui.Box(c).Size(t.Space(20), t.Space(20)).Shrink(0).Radius(t.Radius).Clip().Center().Background(t.Surface)
		border := t.Border
		if root.FileDragOver() {
			border = t.Accent
		}
		frame.Border(1, frame.AnimateColor("border", border, ui.Snappy)).Children(func() {
			if preview != nil {
				ui.Image(c, preview).Fill().Fit(ui.Cover).ContentTransition(ui.ContentTransition{})
			} else {
				ui.Icon(c, icon.Must("image")).Size(16, 16).TextColor(t.TextMuted)
			}
		})
		ui.Column(c).Grow(1).Gap(t.Space(2)).Children(func() {
			if *path == "" {
				kit.Muted(c, o.Empty)
			}
			ui.Row(c).Key("actions").Gap(t.Space(2)).Children(func() {
				if *path == "" {
					pick = plain(c, o.Choose)
					return
				}
				pick = plain(c, o.Replace)
				if plain(c, o.Remove).Clicked() {
					*path = ""
					root.SetChanged()
					sound.Play(sound.Delete)
					kit.Log("ImageField", "removed")
				}
			})
		})
	})
	kit.Clicked(pick, sound.Open)
	return root, pick
}
