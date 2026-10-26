package marketing

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PromoBannerOptions word the strip; Copy takes the code.
type PromoBannerOptions struct{ Copy, Copied string }

// PromoBanner is a promo strip whose code copies; cue Pop.
func PromoBanner(c *ui.Context, text, code string, o PromoBannerOptions) (root, copy ui.Element) {
	def(&o.Copy, "Copy code %s")
	def(&o.Copied, "Copied")
	if strings.TrimSpace(text) == "" || strings.TrimSpace(code) == "" {
		panic(fmt.Sprintf("marketing: PromoBanner %q needs text and a code", code))
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(4)).Radius(t.Radius * 1.5).
		Background(t.Accent).TextColor(t.AccentText).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: -t.Space(2), Fade: 1}, 0).Children(func() {
		ui.Row(c).Grow(1).Basis(t.Space(40)).MinWidth(0).Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			glyph(c, "ticket", 18)
			ui.Text(c, text).FontSize(t.Rem(0.86)).Bold().Grow(1).MinWidth(0).Wrap()
		})
		c.Key("code")
		copy = kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Copy, code)).Row().Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).
			Padding(t.Space(1.5), t.Space(3)).Radius(999).BorderStyle(ui.BorderDashed))
		hover := copy.AnimateTo("hover", unit(copy.Hovered()), ui.Snappy)
		copy.Background(t.AccentText.Alpha(0.12+0.12*hover)).Border(1, t.AccentText.Alpha(0.7))
		done := ui.Local(copy, "until", func() time.Time { return time.Time{} })
		held := ui.Local(copy, "held", func() string { return "" })
		if copy.Clicked() {
			c.WriteClipboard(code)
			*done, *held = c.Now().Add(2*time.Second), code
			sound.Play(sound.Pop)
			kit.Log("PromoBanner", "copied")
		}
		copied := c.Now().Before(*done) && *held == code
		if copied {
			c.After(done.Sub(c.Now()))
		}
		copy.Children(func() {
			shown, mark := code, "copy"
			if copied {
				shown, mark = o.Copied, "check"
			}
			ui.Text(c, shown).FontSize(t.Rem(0.84)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine().ContentTransition(ui.ContentTransition{})
			glyph(c, mark, 14).ContentTransition(ui.ContentTransition{})
		})
	})
	return root, copy
}
