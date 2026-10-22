package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DismissibleAlertOptions tune a DismissibleAlert.
type DismissibleAlertOptions struct {
	// Close names the close button; "Dismiss" when empty.
	Close string
}

var collapse = ui.ElementTransition{Exit: &ui.Motion{Collapse: true}}

// DismissibleAlert is a closable alert that collapses out.
func DismissibleAlert(c *ui.Context, tone Tone, title, body string, open *bool, o DismissibleAlertOptions) ui.Element {
	if open == nil {
		panic("feedback: DismissibleAlert needs open state")
	}
	if o.Close == "" {
		o.Close = "Dismiss"
	}
	root := ui.Column(c).FillWidth()
	return root.Children(func() {
		if !*open {
			return
		}
		c.Key("alert")
		box := alertBox(c, tone, title, body, func() {
			if kit.Clicked(closer(c, o.Close), sound.Close) {
				*open = false
				root.SetChanged()
				kit.Log("DismissibleAlert", "closed")
			}
		})
		kit.Arrive(box, ui.Pose{Y: -6, Fade: 1}, 0).Transition(collapse)
	})
}
