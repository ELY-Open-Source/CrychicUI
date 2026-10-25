package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Attr is one span attribute.
type Attr struct {
	Key, Value string
}

// SpanDetailOptions word the card.
type SpanDetailOptions struct{ Ok, Failed, Attributes, Empty string }

// SpanDetail shows a span's timing and attributes.
func SpanDetail(c *ui.Context, span Span, attrs []Attr, o SpanDetailOptions) ui.Element {
	if span.Name == "" || span.Service == "" || span.Duration < 0 {
		panic(fmt.Sprintf("devtools: SpanDetail %q needs a name, service and duration", span.ID))
	}
	keys := make([]string, len(attrs))
	for i, a := range attrs {
		keys[i] = a.Key
	}
	if len(attrs) > 0 {
		unique("SpanDetail", keys)
	}
	def(&o.Ok, "OK")
	def(&o.Failed, "Error")
	def(&o.Attributes, "Attributes")
	def(&o.Empty, "No attributes")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, span.Name).Font("monospace").FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).Wrap()
				if span.Error {
					kit.Chip(c, o.Failed, t.Danger)
				} else {
					kit.Chip(c, o.Ok, t.Success)
				}
			})
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "server", 14)
					ui.Text(c, span.Service).FontSize(t.Rem(0.78)).Shrink(0)
				})
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "timer", 14)
					kit.Numeric(c, took(span.Duration)).FontSize(t.Rem(0.78)).Shrink(0)
				})
			})
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(2), t.Space(4), t.Space(4), t.Space(4)).Gap(t.Space(0.5))
		list.Children(func() {
			kit.Muted(c, o.Attributes).FontSize(t.Rem(0.72)).Bold().Padding(0, 0, t.Space(1), 0)
			if len(attrs) == 0 {
				empty(c, "tags", o.Empty)
			}
			for i, a := range attrs {
				c.Key(a.Key)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Start).Padding(t.Space(1.25), t.Space(1.5)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					kit.Muted(c, a.Key).Font("monospace").FontSize(t.Rem(0.72)).Width(t.Space(36)).Shrink(0).Wrap()
					ui.Text(c, a.Value).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
	return root
}
