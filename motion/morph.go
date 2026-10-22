package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Morph springs between a dot and a card on click.
func Morph(c *ui.Context, label string, open *bool, content func()) ui.Element {
	if label == "" || open == nil || content == nil {
		panic("motion: Morph needs a label, an open pointer and content")
	}
	t := c.Theme()
	root := ui.ButtonBase(c).Label(label)
	if kit.Clicked(root, sound.Swipe) {
		*open = !*open
		root.SetChanged()
		kit.Log("Morph", "toggled", "open", *open)
	}
	w, h, r := t.Space(12), t.Space(12), t.Space(6)
	if *open {
		w, h, r = t.Space(60), t.Space(30), t.Radius
	}
	a := ui.Bouncy
	root.Width(root.AnimateTo("w", w, a)).Height(root.AnimateTo("h", h, a)).Radius(root.AnimateTo("r", r, a)).
		Clip().Background(t.Accent).TextColor(t.AccentText)
	fade := ui.ElementTransition{Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return kit.Press(root.Checked(*open)).Children(func() {
		if *open {
			ui.Column(c).Padding(t.Space(4)).Transition(fade).Children(content)
		}
	})
}
