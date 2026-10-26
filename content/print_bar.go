package content

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PrintBarOptions word the actions.
type PrintBarOptions struct{ Print, Share, Copied string }

// PrintBar offers print and copy-link actions; cue Tap.
func PrintBar(c *ui.Context, link string, o PrintBarOptions) (root, print, share ui.Element) {
	def(&o.Print, "Print")
	def(&o.Share, "Copy link")
	def(&o.Copied, "Link copied")
	if strings.TrimSpace(link) == "" {
		panic(fmt.Sprintf("content: PrintBar %q needs a link", o.Share))
	}
	t := c.Theme()
	root = ui.Row(c).Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleToolbar)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	pill := func(label string) ui.Element {
		b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(999))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy)).Border(1, t.Border)
		return b
	}
	root.Children(func() {
		c.Key("print")
		print = pill(o.Print)
		print.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "printer", 14)
				ui.Text(c, o.Print).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine()
			})
		})
		if print.Clicked() {
			sound.Play(sound.Tap)
			kit.Log("PrintBar", "print")
		}
		c.Key("share")
		share = pill(o.Share)
		until := ui.Local(share, "until", func() time.Time { return time.Time{} })
		held := ui.Local(share, "held", func() string { return "" })
		if share.Clicked() {
			c.WriteClipboard(link)
			*until, *held = c.Now().Add(2*time.Second), link
			sound.Play(sound.Tap)
			kit.Log("PrintBar", "copied")
		}
		copied := c.Now().Before(*until) && *held == link
		if copied {
			c.After(until.Sub(c.Now()))
		}
		share.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				shown, mark, ink := o.Share, "link", t.Text
				if copied {
					shown, mark, ink = o.Copied, "check", t.Success
				}
				glyph(c, mark, 14).TextColor(ink).ContentTransition(ui.ContentTransition{})
				ui.Text(c, shown).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(36)).SingleLine().
					ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root, print, share
}
