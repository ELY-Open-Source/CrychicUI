package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PageView is one page and its views.
type PageView struct {
	Path  string
	Views int
}

// TopPages shows pages with view bars.
func TopPages(c *ui.Context, pages []PageView) ui.Element {
	if len(pages) == 0 {
		panic("dashboard: TopPages needs pages")
	}
	t := c.Theme()
	maxViews := 0
	for _, page := range pages {
		if page.Path == "" || page.Views < 0 {
			panic("dashboard: TopPages has an invalid page")
		}
		maxViews = max(maxViews, page.Views)
	}
	if maxViews == 0 {
		maxViews = 1
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, page := range pages {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, page.Path).Font("monospace").Grow(1)
					kit.Numeric(c, fmt.Sprint(page.Views)).TextColor(t.TextMuted)
				})
				progress(c, float64(page.Views)/float64(maxViews), t.Accent)
			})
		}
	})
	return root
}
