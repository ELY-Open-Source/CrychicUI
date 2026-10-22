package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InlineFormOptions tune an InlineForm.
type InlineFormOptions struct {
	// Placeholder and Submit are caller copy.
	Placeholder, Submit string
}

// InlineForm is an input and a submit for non-blank text.
func InlineForm(c *ui.Context, label string, value *string, o InlineFormOptions) (root, send ui.Element) {
	if label == "" || value == nil {
		panic("form: InlineForm needs a label and a value")
	}
	def(&o.Submit, "Submit")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
	root.Children(func() {
		edited(root, ui.TextInput(c, value).Label(label).Placeholder(o.Placeholder).Grow(1))
		send = primary(c, o.Submit).Disabled(blank(*value))
	})
	if kit.Clicked(send, sound.Send) {
		kit.Log("InlineForm", "send", "length", len(*value))
	}
	return root, send
}
