package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// IncidentEvent is one timeline entry, oldest first.
type IncidentEvent struct {
	Stage, Text, Time string
}

// IncidentOptions word the card; Meta sits under the title.
type IncidentOptions struct {
	Timeline, Meta, Ongoing, Fixed string
	Resolved                       bool
}

// Incident shows an incident, its status and timeline.
func Incident(c *ui.Context, title string, events []IncidentEvent, o IncidentOptions) ui.Element {
	if title == "" || len(events) == 0 {
		panic(fmt.Sprintf("dashboard: Incident %q needs a title and events", title))
	}
	for _, event := range events {
		if event.Text == "" || event.Time == "" {
			panic(fmt.Sprintf("dashboard: Incident %q event needs text and a time", title))
		}
	}
	def(&o.Timeline, "Timeline")
	def(&o.Ongoing, "Ongoing")
	def(&o.Fixed, "Resolved")
	t := c.Theme()
	ink, glyph, state := t.Danger, "siren", o.Ongoing
	if o.Resolved {
		ink, glyph, state = t.Success, "circle-check", o.Fixed
	}
	root := card(c)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
			Children(func() {
				tile(c, glyph, ink)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, title).Bold().FillWidth().Wrap()
					if o.Meta != "" {
						kit.Muted(c, o.Meta).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
					}
				})
				kit.Chip(c, state, ink)
			})
		list := ui.Column(c).FillWidth().Padding(t.Space(3))
		list.Children(func() {
			kit.Muted(c, o.Timeline).FontSize(t.Rem(0.7)).Bold().Padding(0, 0, t.Space(2), 0)
			for i, event := range events {
				c.Key(i)
				last := i == len(events)-1
				dot := t.TextMuted
				if last {
					dot = ink
				}
				row := ui.Row(c).FillWidth().Gap(t.Space(3))
				kit.Stagger(c, list, row, i).Children(func() {
					rail := ui.Box(c).Width(t.Space(4)).Shrink(0).FillHeight()
					rail.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						cx, d := r.X+r.W/2, t.Space(2.5)
						if !last {
							p.Fill(ui.Rect{X: cx - 0.5, Y: r.Y + t.Space(1) + d, W: 1, H: r.H - d}, t.Border, 0)
						}
						if last {
							p.Fill(ui.Rect{X: cx - d, Y: r.Y + t.Space(1.25) - d/2, W: 2 * d, H: 2 * d}, dot.Alpha(0.18), 999)
						}
						p.Fill(ui.Rect{X: cx - d/2, Y: r.Y + t.Space(1.25), W: d, H: d}, dot, 999)
					}))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Padding(0, 0, t.Space(3), 0).Children(func() {
						if event.Stage != "" {
							ui.Text(c, event.Stage).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
						}
						ui.Text(c, event.Text).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).FillWidth().Wrap()
					})
					kit.Muted(c, event.Time).FontSize(t.Rem(0.75)).Shrink(0)
				})
			}
		})
	})
	return root
}
