package form

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ConditionalField shows content under a switch; its height springs.
func ConditionalField(c *ui.Context, label string, on *bool, content func()) ui.Element {
	if label == "" || on == nil || content == nil {
		panic("form: ConditionalField needs a label, a switch and content")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	spring := ui.Bouncy
	return root.Children(func() {
		if choice.Toggle(c, label, on).Changed() {
			root.SetChanged()
			kit.Log("ConditionalField", "toggled", "on", *on)
		}
		if *on {
			ui.Column(c).Key("more").FillWidth().Gap(t.Space(3)).Clip().
				Transition(ui.ElementTransition{Animation: &spring, Enter: &ui.Motion{Collapse: true}, Exit: &ui.Motion{Collapse: true}}).
				Children(content)
		}
	})
}
