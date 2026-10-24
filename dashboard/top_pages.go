package dashboard

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PageView is one page and its views.
type PageView struct {
	Path  string
	Views int
}

// TopPagesOptions word the card; Views takes the total.
type TopPagesOptions struct{ Title, Views string }

// TopPages ranks pages with a views bar behind each row.
func TopPages(c *ui.Context, pages []PageView, o TopPagesOptions) ui.Element {
	if len(pages) == 0 {
		panic("dashboard: TopPages needs pages")
	}
	names, peak, total := make([]string, len(pages)), 1, 0
	for i, page := range pages {
		if page.Views < 0 {
			panic(fmt.Sprintf("dashboard: TopPages page %q needs views of zero or more", page.Path))
		}
		names[i] = page.Path
		peak = max(peak, page.Views)
		if total > math.MaxInt-page.Views {
			panic("dashboard: TopPages total overflows")
		}
		total += page.Views
	}
	unique("TopPages", names)
	def(&o.Title, "Top pages")
	def(&o.Views, "%s views")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "file-text", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Views, number(c, total)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2))
		list.Children(func() {
			for i, page := range pages {
				c.Key(page.Path)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				w := row.AnimateTo("w", float32(page.Views)/float32(peak), ui.Smooth)
				ink := row.AnimateColor("ink", t.Accent.Alpha(0.08+0.06*unit(row.Hovered())), ui.Snappy)
				row.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, ink, t.Radius)
				}))
				kit.Stagger(c, list, row, i).Children(func() {
					kit.Numeric(c, strconv.Itoa(i+1)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Width(t.Space(4)).Shrink(0)
					ui.Text(c, page.Path).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
					kit.Numeric(c, number(c, page.Views)).Bold().FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
		})
	})
	return root
}
