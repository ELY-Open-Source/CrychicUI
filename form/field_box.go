package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FieldBoxOptions tune a FieldBox.
type FieldBoxOptions struct {
	// Disabled disables every field inside.
	Disabled bool
}

// FieldBox frames fields under a legend; it fades in.
func FieldBox(c *ui.Context, legend string, o FieldBoxOptions, content func()) ui.Element {
	if legend == "" || content == nil {
		panic("form: FieldBox needs a legend and content")
	}
	t := c.Theme()
	box := ui.Fieldset(c, legend, content).Disabled(o.Disabled).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).
		Radius(t.Radius).Border(1, t.Border)
	return kit.Arrive(box, ui.Pose{Fade: 1}, 0)
}
