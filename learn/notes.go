package learn

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Highlight is a quoted passage and note; Tint is 0–3.
type Highlight struct {
	Quote, Note, Where string
	Tint               int
}

// NotesOptions word the list; Count takes a number.
type NotesOptions struct{ Title, Count, One, Empty string }

// Notes lists highlights, each with its note.
func Notes(c *ui.Context, items []Highlight, o NotesOptions) ui.Element {
	def(&o.Title, "Highlights")
	quotes := make([]string, len(items))
	for i, h := range items {
		if strings.TrimSpace(h.Quote) == "" || h.Tint < 0 || h.Tint > 3 {
			panic(fmt.Sprintf("learn: Notes %q highlight %d needs a quote and a tint of 0 to 3", o.Title, i))
		}
		quotes[i] = h.Quote
	}
	unique("Notes", quotes)
	def(&o.Count, "%d highlights")
	def(&o.One, "1 highlight")
	def(&o.Empty, "No highlights yet")
	t := c.Theme()
	inks := [...]ui.Color{t.Warning, t.Success, t.Accent, t.Danger}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		count := fmt.Sprintf(o.Count, len(items))
		if len(items) == 1 {
			count = o.One
		}
		kit.CardHead(c, "highlighter", o.Title, func() { kit.Chip(c, count, t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			if len(items) == 0 {
				ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(5)).TextColor(t.TextMuted).Children(func() {
					glyph(c, "highlighter", 24)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.82)).FillWidth().TextAlign(ui.Center).Wrap()
				})
			}
			for i, h := range items {
				ink := inks[h.Tint]
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Box(c).Width(t.Space(0.75)).FillHeight().Radius(999).Background(ink)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
						ui.Box(c).FillWidth().Padding(t.Space(1), t.Space(2)).Radius(t.Radius * 0.6).Background(ink.Alpha(0.16)).Children(func() {
							ui.Text(c, h.Quote).Font("serif").FontSize(t.Rem(0.86)).FillWidth().Wrap()
						})
						if h.Note != "" {
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Start).Children(func() {
								glyph(c, "sticky-note", 14).TextColor(t.TextMuted)
								ui.Text(c, h.Note).FontSize(t.Rem(0.78)).Grow(1).MinWidth(0).Wrap()
							})
						}
						if h.Where != "" {
							kit.Muted(c, h.Where).FontSize(t.Rem(0.66)).FillWidth().MinWidth(0).SingleLine()
						}
					})
				})
			}
		})
	})
	return root
}
