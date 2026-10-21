package input

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/richtext"
)

// RichEditor is a rich text editor with a styled toolbar.
func RichEditor(c *ui.Context, doc *richtext.Document) ui.Element {
	if doc == nil {
		panic("input: RichEditor needs a document")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		format := doc.Formatting()
		ui.Row(c).Gap(t.Space(1)).Children(func() {
			for _, action := range []struct {
				label, icon string
				on          bool
				run         func()
			}{
				{"Bold", "bold", format.Bold == richtext.FormatOn, doc.ToggleBold},
				{"Italic", "italic", format.Italic == richtext.FormatOn, doc.ToggleItalic},
				{"Underline", "underline", format.Underline == richtext.FormatOn, doc.ToggleUnderline},
			} {
				button := ui.ButtonBase(c).Label(action.label).Padding(t.Space(1)).Radius(t.Radius).Checked(action.on)
				bg := t.Background.Alpha(0)
				if action.on {
					bg = t.Accent.Alpha(0.14)
				}
				kit.Press(button).Background(button.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
					ui.Icon(c, icon.Must(action.icon)).Size(16, 16)
				})
				if kit.Clicked(button, sound.Toggle) {
					action.run()
					root.SetChanged()
					kit.Log("RichEditor", "toggled", "action", action.label)
				}
			}
		})
		editor := richtext.Editor(c, doc).Fill().MinHeight(t.Rem(12))
		if editor.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
		}
	})
}
