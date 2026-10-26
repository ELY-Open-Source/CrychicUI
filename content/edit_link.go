package content

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EditLinkOptions word the link; Date is a layout.
type EditLinkOptions struct{ Label, Updated, Date string }

// EditLink opens the page source; cue Tap.
func EditLink(c *ui.Context, url string, updated time.Time, o EditLinkOptions) (root, edit ui.Element) {
	def(&o.Label, "Edit this page")
	def(&o.Updated, "Updated %s")
	def(&o.Date, "Jan 2, 2006")
	if strings.TrimSpace(url) == "" {
		panic(fmt.Sprintf("content: EditLink %q needs a url", o.Label))
	}
	t := c.Theme()
	root = ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.SpaceBetween).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		edit = kit.Press(ui.ButtonBase(c).Label(o.Label).Role(ui.RoleLink).MinWidth(0).Padding(t.Space(1), t.Space(1.5)).Radius(t.Radius))
		hover := unit(edit.Hovered())
		edit.Background(edit.AnimateColor("bg", t.Accent.Alpha(0.08*hover), ui.Snappy))
		edit.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "pencil", 14).TextColor(t.Accent).Rotate(-12 * edit.AnimateTo("tilt", hover, ui.Snappy))
				ui.Text(c, o.Label).FontSize(t.Rem(0.82)).Bold().TextColor(t.Accent).MinWidth(0).MaxWidth(t.Space(70)).SingleLine()
				glyph(c, "external-link", 12).TextColor(t.Accent.Alpha(0.7))
			})
		})
		if edit.Clicked() {
			c.OpenURL(url)
			sound.Play(sound.Tap)
			kit.Log("EditLink", "opened")
		}
		if !updated.IsZero() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).MinWidth(0).Children(func() {
				glyph(c, "clock", 12).TextColor(t.TextMuted)
				kit.Muted(c, fmt.Sprintf(o.Updated, updated.Format(o.Date))).FontSize(t.Rem(0.74)).MinWidth(0).MaxWidth(t.Space(70)).SingleLine()
			})
		}
	})
	return root, edit
}
