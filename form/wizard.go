package form

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WizardOptions tune a Wizard.
type WizardOptions struct {
	// Back, Next and Finish are button copy.
	Back, Next, Finish string
}

// Wizard pages through steps with Back and Next; pages slide.
func Wizard(c *ui.Context, steps []string, current *int, page func(i int), o WizardOptions) (root, finish ui.Element) {
	if current == nil || page == nil {
		panic("form: Wizard needs a current step and pages")
	}
	def(&o.Back, "Back")
	def(&o.Next, "Next")
	def(&o.Finish, "Finish")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(5))
	dir := ui.Local(root, "dir", func() float32 { return 1 })
	move := func(by int) {
		*current += by
		*dir = float32(by)
		root.SetChanged()
		kit.Log("Wizard", "step", "index", *current)
	}
	root.Children(func() {
		StepHeader(c, steps, *current)
		at, shift := *current, *dir*t.Space(8)
		ui.Box(c).FillWidth().Clip().Children(func() {
			ui.Column(c).Key(at).FillWidth().Gap(t.Space(3)).
				Transition(ui.ElementTransition{Enter: &ui.Motion{X: shift}, Exit: &ui.Motion{X: -shift}}).
				Children(func() { page(at) })
		})
		ui.Row(c).FillWidth().Justify(ui.End).Gap(t.Space(2)).Children(func() {
			if button.Ghost(c, o.Back).Disabled(at == 0).Clicked() {
				move(-1)
			}
			if at < len(steps)-1 {
				if kit.Clicked(primary(c, o.Next), sound.Tap) {
					move(1)
				}
				return
			}
			finish = primary(c, o.Finish)
			if kit.Clicked(finish, sound.Success) {
				kit.Log("Wizard", "finish", "steps", len(steps))
			}
		})
	})
	return root, finish
}
