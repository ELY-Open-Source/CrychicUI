package content

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LinkCardOptions dress the card; Host defaults to the URL's.
type LinkCardOptions struct{ Icon, Host string }

// LinkCard is a titled link with a description; cue Tap.
func LinkCard(c *ui.Context, title, text, link string, o LinkCardOptions) ui.Element {
	def(&o.Icon, "link")
	u, err := url.Parse(link)
	if strings.TrimSpace(title) == "" || err != nil || u.Scheme == "" {
		panic(fmt.Sprintf("content: LinkCard %q needs a title and an absolute url: %v", title, err))
	}
	def(&o.Host, u.Host)
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(title).Role(ui.RoleLink).FillWidth().Padding(t.Space(4)).Radius(t.Radius*1.5).Shadow(0, 1, 3, 0, t.Text.Alpha(0.06)))
	kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	hover := unit(b.Hovered())
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.5*hover), ui.Snappy))
	b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, hover), ui.Snappy))
	shift := b.AnimateTo("nudge", hover, ui.Snappy)
	b.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Children(func() {
			tile(c, o.Icon, t.Accent, t.Space(9))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, title).FontSize(t.Rem(0.9)).Bold().FillWidth().Wrap()
				if text != "" {
					kit.Muted(c, text).FontSize(t.Rem(0.8)).LineHeight(1.5).FillWidth().Wrap()
				}
				ui.Text(c, o.Host).Font("monospace").FontSize(t.Rem(0.7)).TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
			})
			glyph(c, "arrow-up-right", 16).TextColor(t.TextMuted.Mix(t.Accent, hover)).Translate(t.Space(0.75)*shift, -t.Space(0.75)*shift)
		})
	})
	if b.Clicked() {
		c.OpenURL(link)
		sound.Play(sound.Tap)
		kit.Log("LinkCard", "opened")
	}
	return b
}
