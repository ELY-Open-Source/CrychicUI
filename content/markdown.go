package content

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
)

// MarkdownOptions title the view; an empty Title hides the bar.
type MarkdownOptions struct{ Title string }

// Markdown is a styled view of Markdown text.
func Markdown(c *ui.Context, src string, o MarkdownOptions) ui.Element {
	if strings.TrimSpace(src) == "" {
		panic(fmt.Sprintf("content: Markdown %q needs text", o.Title))
	}
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup)
	opts := ui.Local(root, "md", func() *markdown.Options { return &markdown.Options{} })
	return root.Children(func() {
		if o.Title != "" {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(4)).
				Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
				glyph(c, "file-text", 16).TextColor(t.TextMuted)
				ui.Text(c, o.Title).Font("monospace").FontSize(t.Rem(0.78)).Bold().Grow(1).MinWidth(0).SingleLine()
			})
		}
		c.Key("body")
		ui.Column(c).FillWidth().Padding(t.Space(5)).Children(func() {
			markdown.View(c, src, *opts).FillWidth()
		})
	})
}
