package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LabelOptions tune a Label.
type LabelOptions struct {
	// Required shows a danger mark.
	Required bool
}

// Label is medium text with an optional required mark.
func Label(c *ui.Context, text string, o LabelOptions) ui.Element {
	t := c.Theme()
	return ui.Row(c).Gap(t.Space(1)).FontWeight(500).Children(func() {
		ui.Text(c, text).ContentTransition(ui.ContentTransition{})
		if o.Required {
			mark := ui.Text(c, "*").TextColor(t.Danger)
			kit.Arrive(mark, ui.Pose{Scale: 0.5, Fade: 1}, 0)
		}
	})
}
