package feedback

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SuccessState is a drawing check, title and action.
func SuccessState(c *ui.Context, title, text, action string) (root, act ui.Element) {
	if title == "" {
		panic("feedback: SuccessState needs a title")
	}
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Center).MaxWidth(t.Space(100)).AlignItems(ui.Center).Gap(t.Space(2)).
		Padding(t.Space(6)).Role(ui.RoleStatus).Label(title)
	if once(root) {
		sound.Play(sound.Success)
		kit.Log("SuccessState", "shown")
	}
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		d := t.Space(12)
		mark := ui.Box(c).Size(d, d).Radius(999).Background(t.Success.Alpha(0.12)).Margin(0, 0, t.Space(1), 0)
		drawn := mark.Keyframes("draw", grow, ui.KeyframeOptions{Delay: 120 * time.Millisecond})
		mark.Draw(func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.3, r.Y+r.H*0.52).LineTo(r.X+r.W*0.45, r.Y+r.H*0.66).LineTo(r.X+r.W*0.71, r.Y+r.H*0.36)
			p.StrokePath(path.Trim(0, drawn), t.Space(0.75), t.Success)
		})
		ui.Text(c, title).Bold().FontSize(t.Rem(1.15)).TextAlign(ui.Center)
		if text != "" {
			kit.Muted(c, text).TextAlign(ui.Center)
		}
		if action != "" {
			act = button(c, action, true).Margin(t.Space(2), 0, 0, 0)
			if kit.Clicked(act, sound.Tap) {
				kit.Log("SuccessState", "action")
			}
		}
	})
	return root, act
}
