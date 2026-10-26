package geo

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PinOptions tint the pin: accent, good, warning, bad.
type PinOptions struct{ Tint int }

// Pin is a map pin that drops again; cue Drop.
func Pin(c *ui.Context, label string, o PinOptions) ui.Element {
	if strings.TrimSpace(label) == "" || o.Tint < 0 || o.Tint > 3 {
		panic(fmt.Sprintf("geo: Pin %q needs a label and tint 0 to 3", label))
	}
	t := c.Theme()
	ink := [4]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}[o.Tint]
	b := kit.Press(ui.ButtonBase(c).Label(label).Column().AlignItems(ui.Center).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Radius))
	drops := ui.Local(b, "drops", func() int { return 0 })
	if b.Clicked() {
		*drops++
		sound.Play(sound.Drop)
		kit.Log("Pin", "dropped", "drops", *drops)
	}
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	return b.Children(func() {
		// Each drop rekeys the drawing, replaying the fall.
		c.Key(*drops)
		box := ui.Box(c).Size(t.Space(10), t.Space(16)).Shrink(0)
		lift := box.Keyframes("drop", []ui.Keyframe{
			{Value: t.Space(8)},
			{At: 300 * time.Millisecond, Value: 0, Ease: ui.EaseIn},
			{At: 380 * time.Millisecond, Value: t.Space(1.5), Ease: ui.EaseOut},
			{At: 460 * time.Millisecond, Value: 0, Ease: ui.EaseIn},
		})
		box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			size := t.Space(4)
			tip := r.Y + r.H - size*0.5
			shadow := float32(1) - min(1, lift/t.Space(8))*0.7
			var under ui.Path
			p.FillPath(under.Circle(r.X+r.W/2, tip, size*0.45*shadow), t.Text.Alpha(0.18*shadow))
			pin(p, r.X+r.W/2, tip-lift, size, ink, t.Background)
		}))
		ui.Text(c, label).FontSize(t.Rem(0.72)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine()
	})
}
