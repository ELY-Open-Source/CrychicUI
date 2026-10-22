package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WizardNavOptions tune a WizardNav.
type WizardNavOptions struct {
	Back, Next, Finish string
}

// WizardNav is back, step dots and next or finish.
func WizardNav(c *ui.Context, label string, step *int, n int, o WizardNavOptions) (root, finish ui.Element) {
	if label == "" || n < 2 || step == nil || !valid(*step, n) {
		panic(fmt.Sprintf("navigation: WizardNav %q needs two steps and a current one", label))
	}
	def(&o.Back, "Back")
	def(&o.Next, "Next")
	def(&o.Finish, "Finish")
	t := c.Theme()
	root = ui.Row(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(4)).Label(label)
	shown := *step
	root.Children(func() {
		back := button.Ghost(c.Key("back"), o.Back).Disabled(shown == 0)
		if back.Clicked() {
			*step = shown - 1
			root.SetChanged()
			kit.Log("WizardNav", "moved", "step", *step)
		}
		ui.Row(c).Grow(1).Gap(t.Space(1.5)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			for i := range n {
				c.Key(i)
				pip(c, ui.Box(c), i == shown)
			}
		})
		if shown == n-1 {
			finish = button.Solid(c.Key("finish"), o.Finish)
			if finish.Clicked() {
				kit.Log("WizardNav", "finished")
			}
			return
		}
		if button.Solid(c.Key("next"), o.Next).Clicked() {
			*step = shown + 1
			root.SetChanged()
			kit.Log("WizardNav", "moved", "step", *step)
		}
	})
	return root, finish
}
