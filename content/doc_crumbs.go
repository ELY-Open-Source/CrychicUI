package content

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DocCrumbs is a doc page's path, current part last.
func DocCrumbs(c *ui.Context, path []string) ui.Element {
	if len(path) == 0 {
		panic("content: DocCrumbs needs a path")
	}
	for _, p := range path {
		if strings.TrimSpace(p) == "" {
			panic("content: DocCrumbs has a blank part")
		}
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	return root.Children(func() {
		glyph(c, "book-open", 14).TextColor(t.TextMuted)
		for i, p := range path {
			c.Key(i)
			last := i == len(path)-1
			part := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).MinWidth(0)
			kit.Stagger(c, root, part, i).Children(func() {
				glyph(c, "chevron-right", 12).TextColor(t.TextMuted.Alpha(0.6))
				txt := ui.Text(c, p).FontSize(t.Rem(0.78)).MinWidth(0).MaxWidth(t.Space(40)).SingleLine()
				if last {
					txt.Bold().TextColor(t.Text)
				} else {
					txt.TextColor(t.TextMuted)
				}
			})
		}
	})
}
