package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CTAOptions word the band; Note shows under the action.
type CTAOptions struct{ Note string }

// CTA is a call-to-action band with one action; cue Tap.
func CTA(c *ui.Context, title, text, act string, o CTAOptions) (root, button ui.Element) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(act) == "" {
		panic(fmt.Sprintf("marketing: CTA %q needs a title and an action", title))
	}
	t := c.Theme()
	root = card(c).Gap(t.Space(3)).Padding(t.Space(8), t.Space(5)).AlignItems(ui.Center).Clip().Border(0, t.Border).
		Background(t.Accent).TextColor(t.AccentText).Role(ui.RoleGroup)
	drift := root.Keyframes("drift", []ui.Keyframe{{Value: 0}, {At: 1200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Draw(kit.Inside(t.Radius*1.5, func(p *ui.Painter, r ui.Rect) {
		p.FillPath(circle(r.X+r.W*0.92, r.Y+r.H*(0.1+0.1*drift), r.H*0.55), t.AccentText.Alpha(0.08))
		p.FillPath(circle(r.X+r.W*0.05, r.Y+r.H*(1-0.1*drift), r.H*0.4), t.AccentText.Alpha(0.06))
	}))
	root.Children(func() {
		kit.Arrive(ui.Text(c, title).FontSize(t.Rem(1.5)).Bold().FillWidth().TextAlign(ui.Center).Wrap(), ui.Pose{Y: t.Space(3), Fade: 1}, 0)
		if text != "" {
			ui.Text(c, text).FontSize(t.Rem(0.9)).TextColor(t.AccentText.Alpha(0.85)).FillWidth().MaxWidth(t.Space(110)).TextAlign(ui.Center).Wrap()
		}
		c.Key("act")
		button = kit.Press(ui.ButtonBase(c).Label(act).Shrink(0).Padding(t.Space(2.5), t.Space(5)).Radius(999))
		hover := button.AnimateTo("hover", unit(button.Hovered()), ui.Snappy)
		button.Background(t.AccentText.Mix(t.SurfaceHover, 0.3*hover)).TextColor(t.Accent)
		kit.Arrive(button, ui.Pose{Y: t.Space(2), Fade: 1}, 150*time.Millisecond).Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, act).FontSize(t.Rem(0.88)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
				glyph(c, "arrow-right", 16).Translate(t.Space(0.75)*hover, 0)
			})
		})
		if o.Note != "" {
			ui.Text(c, o.Note).FontSize(t.Rem(0.74)).TextColor(t.AccentText.Alpha(0.75)).FillWidth().TextAlign(ui.Center).Wrap()
		}
	})
	if button.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("CTA", "action")
	}
	return root, button
}
