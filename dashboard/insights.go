package dashboard

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Insight is one generated insight; Metric is a signed figure.
type Insight struct {
	Title, Body, Icon, Metric string
}

// InsightsOptions word the card.
type InsightsOptions struct{ Title, Empty string }

// Insights lists generated insight cards, staggered in.
func Insights(c *ui.Context, cards []Insight, o InsightsOptions) ui.Element {
	names := make([]string, len(cards))
	for i, card := range cards {
		if card.Body == "" {
			panic(fmt.Sprintf("dashboard: Insights card %q needs a body", card.Title))
		}
		names[i] = card.Title
	}
	unique("Insights", names)
	def(&o.Title, "Insights")
	def(&o.Empty, "No insights yet")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "sparkles", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		list.Children(func() {
			if len(cards) == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).
					Children(func() {
						ui.Icon(c, icon.Must("lightbulb")).Size(20, 20).TextColor(t.TextMuted)
						kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
					})
			}
			for i, card := range cards {
				c.Key(card.Title)
				glyph := card.Icon
				if glyph == "" {
					glyph = "lightbulb"
				}
				panel := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Radius(t.Radius)
				panel.Background(panel.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(panel.Hovered())), ui.Snappy))
				kit.Stagger(c, list, panel, i).Children(func() {
					tile(c, glyph, palette(t, i))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, card.Title).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
							switch {
							case card.Metric == "":
							case strings.HasPrefix(card.Metric, "+") || falling(card.Metric):
								deltaChip(c, card.Metric)
							default:
								kit.Chip(c, card.Metric, t.Accent)
							}
						})
						ui.Text(c, card.Body).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).FillWidth().Wrap()
					})
				})
			}
		})
	})
	return root
}
