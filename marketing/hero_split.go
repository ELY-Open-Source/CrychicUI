package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HeroSplitOptions word the hero; Points list under the text.
type HeroSplitOptions struct {
	Eyebrow string
	Points  []string
}

// HeroSplit is a hero with media beside; cue Tap.
func HeroSplit(c *ui.Context, title, text, act string, media *ui.Bitmap, o HeroSplitOptions) (root, button ui.Element) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(act) == "" || media == nil {
		panic(fmt.Sprintf("marketing: HeroSplit %q needs a title, an action and media", title))
	}
	unique("HeroSplit", o.Points)
	t := c.Theme()
	root = card(c).Clip()
	root.Children(func() {
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(6)).AlignItems(ui.Center).Padding(t.Space(6)).Children(func() {
			col := ui.Column(c).Grow(1).Basis(t.Space(48)).MinWidth(0).Gap(t.Space(3))
			w := col.LayoutBounds().W
			kit.Settle(c, col, []float32{w}, func() []float32 { return []float32{col.LayoutBounds().W} })
			room := t.Space(60)
			if w > 0 {
				room = max(t.Space(10), min(room, w-2*t.Space(4)-16-t.Space(1.5)))
			}
			col.Children(func() {
				if o.Eyebrow != "" {
					kit.Arrive(ui.Text(c, o.Eyebrow).FontSize(t.Rem(0.74)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine(),
						ui.Pose{Y: t.Space(2), Fade: 1}, 0)
				}
				kit.Arrive(ui.Text(c, title).FontSize(t.Rem(1.5)).Bold().FillWidth().Wrap(), ui.Pose{Y: t.Space(3), Fade: 1}, 60*time.Millisecond)
				if text != "" {
					kit.Muted(c, text).FontSize(t.Rem(0.88)).FillWidth().Wrap()
				}
				list := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				list.Children(func() {
					for i, p := range o.Points {
						r := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
						kit.Stagger(c, list, r, i).Children(func() {
							ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999).Background(t.Success.Alpha(0.14)).
								TextColor(t.Success).Children(func() { glyph(c, "check", 12) })
							ui.Text(c, p).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).Wrap()
						})
					}
				})
				c.Key("act")
				ui.Row(c).Padding(t.Space(1), 0, 0, 0).Children(func() { button = action(c, act, "arrow-right", true, room) })
			})
			frame := ui.Box(c).Grow(1).Basis(t.Space(48)).MinWidth(0).AspectRatio(4.0 / 3).Radius(t.Radius * 1.5).Clip().
				Background(t.SurfaceHover)
			kit.Arrive(frame, ui.Pose{X: t.Space(4), Scale: 0.96, Fade: 1}, 120*time.Millisecond).Children(func() {
				ui.Image(c, media).Fill().Fit(ui.Cover)
			})
		})
	})
	if button.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("HeroSplit", "action")
	}
	return root, button
}
