package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AnnouncePill is a tagged news link; cue Tap.
func AnnouncePill(c *ui.Context, tag, text string) ui.Element {
	if strings.TrimSpace(tag) == "" || strings.TrimSpace(text) == "" {
		panic(fmt.Sprintf("marketing: AnnouncePill %q needs a tag and text", text))
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(text).Shrink(0).Padding(t.Space(1), t.Space(3), t.Space(1), t.Space(1)).Radius(999).
		Border(1, t.Border))
	hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
	b.Background(t.Background.Mix(t.SurfaceHover, hover))
	sheen := b.Keyframes("sheen", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 0},
		{At: 1300 * time.Millisecond, Value: 1, Ease: ui.EaseInOut}})
	b.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		if sheen > 0 && sheen < 1 {
			w := r.H * 2
			p.Fill(ui.Rect{X: r.X - w + (r.W+w)*sheen, Y: r.Y, W: w, H: r.H}, t.Accent.Alpha(0.12), 0)
		}
	}))
	kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, tag).FontSize(t.Rem(0.7)).Bold().TextColor(t.AccentText).Shrink(0).Padding(t.Space(0.5), t.Space(2)).
				Radius(999).Background(t.Accent)
			ui.Text(c, text).FontSize(t.Rem(0.8)).MinWidth(0).MaxWidth(t.Space(80)).SingleLine()
			glyph(c, "arrow-right", 14).TextColor(t.TextMuted).Translate(t.Space(0.75)*hover, 0)
		})
	})
	if b.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("AnnouncePill", "opened")
	}
	return b
}
