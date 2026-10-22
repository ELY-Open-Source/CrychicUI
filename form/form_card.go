package form

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FormCardOptions tune a FormCard.
type FormCardOptions struct {
	// Description, Submit and Cancel are caller copy.
	Description, Submit, Cancel string
}

// FormCard holds titled fields on a card; it rises in.
func FormCard(c *ui.Context, title string, o FormCardOptions, content func()) (root, submit, cancel ui.Element) {
	if title == "" || content == nil {
		panic("form: FormCard needs a title and content")
	}
	def(&o.Submit, "Save")
	def(&o.Cancel, "Cancel")
	t := c.Theme()
	root = kit.Arrive(kit.Card(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(6)), ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			ui.Text(c, title).Bold().FontSize(t.Rem(1.25))
			if o.Description != "" {
				kit.Muted(c, o.Description)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(4)).Children(content)
		ui.Row(c).FillWidth().Justify(ui.End).Gap(t.Space(2)).Children(func() {
			cancel = button.Ghost(c, o.Cancel)
			submit = button.Solid(c, o.Submit)
		})
	})
	return root, submit, cancel
}
