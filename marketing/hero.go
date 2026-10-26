package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HeroOptions word the hero; Eyebrow shows above the title.
type HeroOptions struct{ Eyebrow string }

// Hero is a headline with two actions; cue Tap.
func Hero(c *ui.Context, title, text, primary, secondary string, o HeroOptions) (root, first, second ui.Element) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(primary) == "" || strings.TrimSpace(secondary) == "" || primary == secondary {
		panic(fmt.Sprintf("marketing: Hero %q needs a title and two different actions", title))
	}
	t := c.Theme()
	root = card(c).Gap(t.Space(4)).Padding(t.Space(10), t.Space(6)).AlignItems(ui.Center).Background(t.Background.Mix(t.Accent, 0.05))
	rise := func(e ui.Element, i int) ui.Element {
		return kit.Arrive(e, ui.Pose{Y: t.Space(4), Fade: 1}, time.Duration(i)*90*time.Millisecond)
	}
	root.Children(func() {
		if o.Eyebrow != "" {
			pill := ui.Row(c).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(3)).Radius(999).
				Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
			rise(pill, 0).Children(func() {
				glyph(c, "sparkles", 14)
				ui.Text(c, o.Eyebrow).FontSize(t.Rem(0.74)).Bold().Shrink(0).MaxWidth(t.Space(70)).SingleLine()
			})
		}
		rise(ui.Text(c, title).FontSize(t.Rem(1.9)).Bold().FillWidth().TextAlign(ui.Center).Wrap(), 1)
		if text != "" {
			rise(kit.Muted(c, text).FontSize(t.Rem(0.95)).FillWidth().MaxWidth(t.Space(120)).TextAlign(ui.Center).Wrap(), 2)
		}
		c.Key("actions")
		rise(ui.Row(c).Wrap().Gap(t.Space(2)).Justify(ui.Center).Padding(t.Space(2), 0, 0, 0), 3).Children(func() {
			first = action(c, primary, "arrow-right", true, t.Space(60))
			second = action(c, secondary, "", false, t.Space(60))
		})
	})
	if first.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("Hero", "action", "primary", true)
	}
	if second.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("Hero", "action", "primary", false)
	}
	return root, first, second
}
