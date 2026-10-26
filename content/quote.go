package content

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QuoteOptions add the speaker's role and the work cited.
type QuoteOptions struct{ Role, Cite string }

// Quote is a pull quote with its source.
func Quote(c *ui.Context, text, source string, o QuoteOptions) ui.Element {
	if strings.TrimSpace(text) == "" || strings.TrimSpace(source) == "" {
		panic(fmt.Sprintf("content: Quote by %q needs text and a source", source))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(5)).
		BorderWidth(0, 0, 0, t.Space(1)).BorderColor(t.Accent.Alpha(0.5)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		kit.Arrive(glyph(c, "quote", 24).TextColor(t.Accent), ui.Pose{Scale: 0.5, Rotate: -15, Fade: 1}, 100)
		ui.Text(c, text).FontSize(t.Rem(1.1)).LineHeight(1.5).FillWidth().Wrap()
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Width(t.Space(5)).Height(1).Shrink(0).Background(t.TextMuted)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, source).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
				meta(c, o.Role, o.Cite)
			})
		})
	})
}
