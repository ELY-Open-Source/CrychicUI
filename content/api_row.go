package content

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// APIRowOptions mark a deprecated route.
type APIRowOptions struct {
	Deprecated bool
	Label      string
}

// APIRow is a method, a path and a description.
func APIRow(c *ui.Context, method, path, text string, o APIRowOptions) ui.Element {
	def(&o.Label, "Deprecated")
	method = strings.ToUpper(strings.TrimSpace(method))
	if method == "" || !strings.HasPrefix(path, "/") {
		panic(fmt.Sprintf("content: APIRow %q needs a method and a path starting with /", path))
	}
	t := c.Theme()
	ink := map[string]ui.Color{"GET": t.Success, "POST": t.Accent, "PUT": t.Warning, "PATCH": t.Warning, "DELETE": t.Danger}[method]
	if ink.A == 0 {
		ink = t.TextMuted
	}
	root := card(c).Row().Gap(t.Space(3)).Padding(t.Space(3), t.Space(4)).AlignItems(ui.Start).Role(ui.RoleGroup)
	return root.Children(func() {
		ui.Text(c, method).Font("monospace").FontSize(t.Rem(0.7)).Bold().TextColor(ink).Width(t.Space(16)).Shrink(0).TextAlign(ui.Center).
			Padding(t.Space(0.75), t.Space(1)).Radius(t.Radius).Background(ink.Alpha(0.12)).SingleLine()
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				var spans []ui.Span
				for _, seg := range strings.Split(path[1:], "/") {
					col := t.Text
					if strings.HasPrefix(seg, "{") || strings.HasPrefix(seg, ":") {
						col = t.Accent
					}
					spans = append(spans, ui.Span{Text: "/", Color: t.TextMuted}, ui.Span{Text: seg, Color: col, Weight: 600})
				}
				p := ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.84)).Grow(1).MinWidth(0).Wrap()
				if o.Deprecated {
					p.Strikethrough().Opacity(0.7)
					ui.Text(c, o.Label).FontSize(t.Rem(0.66)).Bold().TextColor(t.Warning).Padding(t.Space(0.25), t.Space(1.5)).Radius(999).
						Background(t.Warning.Alpha(0.12)).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
				}
			})
			if text != "" {
				kit.Muted(c, text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
			}
		})
	})
}
