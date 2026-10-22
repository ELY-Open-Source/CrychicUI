package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Snackbar is a dark action bar that slides up.
func Snackbar(c *ui.Context, message, action string, open *bool) (root, act ui.Element) {
	if message == "" || open == nil {
		panic("feedback: Snackbar needs a message and open state")
	}
	t := c.Theme()
	fill, ink := t.Inverse, t.InverseText
	if fill.A == 0 {
		fill = t.Text
	}
	if ink.A == 0 {
		ink = t.Background
	}
	root = ui.Column(c).AlignItems(ui.Center).FillWidth()
	root.Children(func() {
		if !*open {
			return
		}
		c.Key("bar")
		bar := ui.Row(c).Gap(t.Space(4)).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(4)).Radius(t.Radius).
			Background(fill).TextColor(ink).Shadow(0, 6, 20, 0, kit.Shade(0.2)).Role(ui.RoleStatus).Label(message)
		if once(bar) {
			sound.Play(sound.Notify)
			kit.Log("Snackbar", "shown")
		}
		kit.Arrive(bar, ui.Pose{Y: t.Space(6), Fade: 1}, 0).Transition(ui.ElementTransition{Exit: &ui.Motion{Y: t.Space(6)}})
		bar.Children(func() {
			ui.Text(c, message)
			if action == "" {
				return
			}
			act = ui.ButtonBase(c).Label(action).Padding(t.Space(1), t.Space(2)).Radius(t.Radius)
			bg := ui.Color{}
			if act.Hovered() {
				bg = ink.Alpha(0.12)
			}
			act.Background(act.AnimateColor("bg", bg, quick)).TextColor(t.Accent.Mix(ink, 0.35))
			kit.Press(act).Children(func() { ui.Text(c, action).Bold() })
			if kit.Clicked(act, sound.Tap) {
				*open = false
				root.SetChanged()
				kit.Log("Snackbar", "action")
			}
		})
	})
	return root, act
}
