package content

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DocPage is a searchable page in a section.
type DocPage struct{ Section, Title, Text string }

// DocSearchOptions word the search.
type DocSearchOptions struct{ Placeholder, Empty string }

// DocSearch filters pages by query, grouped by section; cue Select.
func DocSearch(c *ui.Context, query *string, pages []DocPage, picked *string, o DocSearchOptions) ui.Element {
	def(&o.Placeholder, "Search docs")
	def(&o.Empty, "No pages match")
	titles := make([]string, len(pages))
	for i, p := range pages {
		if strings.TrimSpace(p.Section) == "" {
			panic(fmt.Sprintf("content: DocSearch page %q needs a section", p.Title))
		}
		titles[i] = p.Title
	}
	unique("DocSearch", titles)
	if query == nil || picked == nil {
		panic("content: DocSearch needs a query and a picked page")
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		frame := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(999).Background(t.Background)
		frame.Border(1, frame.AnimateColor("ring", t.Border.Mix(t.Accent, unit(frame.FocusWithin())), ui.Snappy))
		frame.Children(func() {
			glyph(c, "search", 16).TextColor(t.TextMuted)
			in := ui.TextInputBase(c, query).Label(o.Placeholder).Placeholder(o.Placeholder).FontSize(t.Rem(0.86)).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				kit.Log("DocSearch", "query", "runes", utf8.RuneCountInString(*query))
			}
		})
		q := strings.TrimSpace(*query)
		var hits []DocPage
		var sections []string
		for _, p := range pages {
			if lo, _ := kit.IndexFold(p.Title, q); q != "" && lo < 0 {
				if lo, _ := kit.IndexFold(p.Text, q); lo < 0 {
					continue
				}
			}
			hits = append(hits, p)
			if !slices.Contains(sections, p.Section) {
				sections = append(sections, p.Section)
			}
		}
		if len(hits) == 0 {
			ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(6), 0).Children(func() {
				kit.Arrive(glyph(c, "search-x", 24).TextColor(t.TextMuted), ui.Pose{Scale: 0.7, Fade: 1}, 0)
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).FillWidth().TextAlign(ui.Center).Wrap()
			})
			return
		}
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleList)
		list.Children(func() {
			for _, sec := range sections {
				c.Key(strconv.Quote(sec))
				ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
					kit.Muted(c, sec).FontSize(t.Rem(0.7)).Bold().FillWidth().MinWidth(0).SingleLine().Padding(0, t.Space(2))
					for _, p := range hits {
						if p.Section != sec {
							continue
						}
						c.Key(strconv.Quote(p.Title))
						on := p.Title == *picked
						b := kit.Press(ui.ButtonBase(c).Label(p.Title).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Checked(on))
						b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
						b.Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
								glyph(c, "file-text", 16).TextColor(t.TextMuted.Mix(t.Accent, unit(on)))
								ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
									lo, hi := kit.IndexFold(p.Title, q)
									ui.RichText(c, kit.Marked(t, p.Title, lo, hi)...).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
									if p.Text != "" {
										lo, hi := kit.IndexFold(p.Text, q)
										ui.RichText(c, kit.Marked(t, p.Text, lo, hi)...).FontSize(t.Rem(0.74)).TextColor(t.TextMuted).
											FillWidth().MinWidth(0).SingleLine()
									}
								})
								glyph(c, "corner-down-left", 14).TextColor(t.TextMuted).Opacity(unit(b.Hovered() || on))
							})
						})
						if b.Clicked() && !on {
							*picked = p.Title
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("DocSearch", "picked")
						}
					}
				})
			}
		})
	})
	return root
}
