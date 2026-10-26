package content

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
)

// EmbedOptions dress the frame's bar.
type EmbedOptions struct{ Icon, Host string }

// Embed frames embedded content under a titled bar.
func Embed(c *ui.Context, title string, o EmbedOptions, body func()) ui.Element {
	def(&o.Icon, "app-window")
	if strings.TrimSpace(title) == "" || body == nil {
		panic(fmt.Sprintf("content: Embed %q needs a title and a body", title))
	}
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup).Label(title)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
			Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			glyph(c, o.Icon, 14).TextColor(t.Accent)
			ui.Text(c, title).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
			if o.Host != "" {
				ui.Text(c, o.Host).Font("monospace").FontSize(t.Rem(0.66)).TextColor(t.TextMuted).Padding(t.Space(0.25), t.Space(1.5)).
					Radius(999).Background(t.Background).Shrink(0).MaxWidth(t.Space(36)).SingleLine()
			}
		})
		ui.Column(c).FillWidth().Padding(t.Space(4)).Children(body)
	})
}
