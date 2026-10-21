package input

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NoteInput is a borderless text area with a fading placeholder.
func NoteInput(c *ui.Context, value *string, placeholder string) ui.Element {
	if value == nil {
		panic("input: NoteInput needs a value")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().MinHeight(t.Rem(10))
	return root.Children(func() {
		in := ui.TextAreaBase(c, value).Fill().Padding(t.Space(3))
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
		}
		if *value == "" {
			placeholderText := ui.Text(c, placeholder).Absolute().Left(t.Space(3)).Top(t.Space(3)).TextColor(t.TextMuted).PassThrough()
			alpha := 1.0
			if in.Focused() {
				alpha = 0.35
			}
			placeholderText.Opacity(placeholderText.AnimateTo("fade", float32(alpha), ui.Snappy))
		}
	})
}
