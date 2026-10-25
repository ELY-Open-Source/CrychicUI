package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Change is one changelog entry; Who, From, To optional.
type Change struct {
	Who, Text, Time, From, To string
}

// ChangesOptions word the card.
type ChangesOptions struct{ Title, Empty string }

// Changes shows a task changelog on a rail.
func Changes(c *ui.Context, entries []Change, o ChangesOptions) ui.Element {
	def(&o.Title, "History")
	def(&o.Empty, "No changes yet")
	for i, entry := range entries {
		if entry.Text == "" || entry.Time == "" || (entry.From == "") != (entry.To == "") {
			panic(fmt.Sprintf("project: Changes %q entry %d needs text, time, and both or neither of from and to", o.Title, i))
		}
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "history", o.Title, func() { kit.Chip(c, fmt.Sprint(len(entries)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Padding(t.Space(3))
		list.Children(func() {
			if len(entries) == 0 {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "history", 16)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
			for i, entry := range entries {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Stretch)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Width(t.Space(7)).Shrink(0).AlignItems(ui.Center).Children(func() {
						if entry.Who != "" {
							display.Avatar(c, entry.Who, nil).Size(t.Space(7), t.Space(7)).FontSize(t.Rem(0.55))
						} else {
							ui.Box(c).Size(t.Space(7), t.Space(7)).Center().Children(func() {
								ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Accent)
							})
						}
						if i < len(entries)-1 {
							ui.Box(c).Width(1).Grow(1).MinHeight(t.Space(2)).Background(t.Border)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(1), 0, t.Space(3), 0).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Start).Children(func() {
							ui.Row(c).Grow(1).MinWidth(0).Wrap().Gap(t.Space(1)).Children(func() {
								if entry.Who != "" {
									ui.Text(c, entry.Who).FontSize(t.Rem(0.8)).Bold().Shrink(0)
								}
								ui.Text(c, entry.Text).FontSize(t.Rem(0.8)).TextColor(t.Text.Mix(t.TextMuted, 0.4)).MinWidth(0).Wrap()
							})
							kit.Muted(c, entry.Time).FontSize(t.Rem(0.7)).Shrink(0)
						})
						if entry.From != "" {
							ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								tag(c, entry.From, t.TextMuted)
								glyph(c, "arrow-right", 12).TextColor(t.TextMuted)
								tag(c, entry.To, t.Accent)
							})
						}
					})
				})
			}
		})
	})
	return root
}
