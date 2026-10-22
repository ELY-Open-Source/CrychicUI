package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FormSectionOptions tune a FormSection.
type FormSectionOptions struct {
	// Description is caller copy under the title.
	Description string
}

// FormSection groups fields under a title; it fades in.
func FormSection(c *ui.Context, title string, o FormSectionOptions, content func()) ui.Element {
	if title == "" || content == nil {
		panic("form: FormSection needs a title and content")
	}
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).FillWidth().Gap(t.Space(3)), ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		ui.Column(c).Gap(t.Space(1)).Children(func() {
			ui.Text(c, title).FontWeight(600)
			if o.Description != "" {
				kit.Muted(c, o.Description)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(4)).Children(content)
	})
}
