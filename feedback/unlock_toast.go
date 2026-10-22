package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

var shine = []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}}

// UnlockToast is an achievement toast with a shine passing once.
func UnlockToast(c *ui.Context, title, text string, open *bool) ui.Element {
	if title == "" || open == nil {
		panic("feedback: UnlockToast needs a title and open state")
	}
	t := c.Theme()
	root := ui.Column(c).AlignItems(ui.Start)
	return root.Children(func() {
		if !*open {
			return
		}
		c.Key("toast")
		card := toastFace(c, ui.Row(c)).AlignItems(ui.Center).Clip().Role(ui.RoleStatus).Label(title)
		if once(card) {
			sound.Play(sound.Success)
			kit.Log("UnlockToast", "shown")
		}
		kit.Arrive(card, ui.Pose{Scale: 0.9, Fade: 1}, 0).Transition(ui.ElementTransition{Exit: &ui.Motion{Scale: 0.95}})
		k := card.Keyframes("shine", shine, ui.KeyframeOptions{Delay: 200 * time.Millisecond})
		card.DrawOver(kit.Inside(t.Radius+2, func(p *ui.Painter, r ui.Rect) {
			if k <= 0 || k >= 1 {
				return
			}
			w := r.W * 0.35
			x := r.X - w + (r.W+w)*k
			band := ui.LinearGradient{Angle: 90, Stops: []ui.GradientStop{{Color: t.Warning.Alpha(0), At: 0}, {Color: t.Warning.Alpha(0.25), At: 0.5}, {Color: t.Warning.Alpha(0), At: 1}}}
			p.FillGradient(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, band, 0)
		}))
		card.Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(t.Warning.Alpha(0.16)).
				TextColor(t.Warning).Children(func() { ui.Icon(c, icon.Must("trophy")).Size(t.Space(5), t.Space(5)) })
			ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, title).Bold()
				if text != "" {
					kit.Muted(c, text)
				}
			})
		})
	})
}
