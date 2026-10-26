package content

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DocSection is a titled group of pages.
type DocSection struct {
	Title, Icon string
	Pages       []string
}

// DocSidebar is a doc tree; folds report too; cue Select.
func DocSidebar(c *ui.Context, sections []DocSection, current *string) ui.Element {
	titles, pages := make([]string, len(sections)), []string{}
	for i, s := range sections {
		if s.Icon == "" || len(s.Pages) == 0 {
			panic(fmt.Sprintf("content: DocSidebar section %q needs an icon and pages", s.Title))
		}
		titles[i] = s.Title
		pages = append(pages, s.Pages...)
	}
	unique("DocSidebar", titles)
	unique("DocSidebar", pages)
	if len(sections) == 0 || current == nil || !slices.Contains(pages, *current) {
		panic("content: DocSidebar needs sections and a current page among them")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleTree)
	kit.Arrive(root, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
	open := ui.Local(root, "open", func() map[string]bool {
		m := map[string]bool{}
		for _, s := range sections {
			m[s.Title] = slices.Contains(s.Pages, *current)
		}
		return m
	})
	root.Children(func() {
		for _, s := range sections {
			c.Key(strconv.Quote(s.Title))
			ui.Column(c).FillWidth().Children(func() {
				shown := (*open)[s.Title]
				h := kit.Press(ui.ButtonBase(c).Label(s.Title).FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Expanded(shown))
				h.Background(h.AnimateColor("bg", t.SurfaceHover.Alpha(unit(h.Hovered())), ui.Snappy))
				turn := h.AnimateTo("turn", unit(shown), ui.Snappy)
				h.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						glyph(c, s.Icon, 16).TextColor(t.TextMuted)
						ui.Text(c, s.Title).FontSize(t.Rem(0.84)).Bold().Grow(1).MinWidth(0).SingleLine()
						glyph(c, "chevron-right", 14).TextColor(t.TextMuted).Rotate(90 * turn)
					})
				})
				if h.Clicked() {
					(*open)[s.Title] = !shown
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("DocSidebar", "section", "open", !shown)
				}
				layout.Collapse(c, shown, func() {
					list := ui.Column(c).FillWidth().Margin(0, 0, 0, t.Space(4)).Padding(t.Space(0.5), 0, t.Space(0.5), t.Space(2)).
						BorderWidth(0, 0, 0, 1).BorderColor(t.Border)
					list.Children(func() {
						for i, page := range s.Pages {
							c.Key(strconv.Quote(page))
							on := page == *current
							b := ui.ButtonBase(c).Label(page).FillWidth().Padding(t.Space(1.25), t.Space(2)).Radius(t.Radius).
								AlignItems(ui.Center).Justify(ui.Start).Checked(on)
							b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
							b.Children(func() {
								txt := ui.Text(c, page).FontSize(t.Rem(0.8)).TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(on)), ui.Snappy)).
									Grow(1).MinWidth(0).SingleLine()
								if on {
									txt.Bold()
								}
							})
							if b.Clicked() && !on {
								*current = page
								root.SetChanged()
								sound.Play(sound.Select)
								kit.Log("DocSidebar", "page", "index", i)
							}
						}
					})
				})
			})
		}
	})
	return root
}
