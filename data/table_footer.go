package data

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TableFooterOptions tune a TableFooter.
type TableFooterOptions struct {
	// Sizes offered; nil gives 10, 25, 50.
	Sizes []int
	// Range takes first, last and total rows.
	Rows, Range, Prev, Next string
}

// TableFooter picks page size and pages from 0.
func TableFooter(c *ui.Context, total int, page, size *int, o TableFooterOptions) ui.Element {
	if o.Sizes == nil {
		o.Sizes = []int{10, 25, 50}
	}
	if total < 0 || page == nil || size == nil || !slices.Contains(o.Sizes, *size) {
		panic("data: TableFooter needs a total, a page and an offered size")
	}
	for _, s := range o.Sizes {
		if s <= 0 {
			panic(fmt.Sprintf("data: TableFooter size %d is not positive", s))
		}
	}
	pages := max(1, (total+*size-1) / *size)
	if !valid(*page, pages) {
		panic(fmt.Sprintf("data: TableFooter page %d is not among %d", *page, pages))
	}
	def(&o.Rows, "Rows per page")
	def(&o.Range, "%d–%d of %d")
	def(&o.Prev, "Previous page")
	def(&o.Next, "Next page")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Justify(ui.End).Padding(t.Space(2), t.Space(3))
	turn := func(to int) {
		*page = to
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("TableFooter", "page", "index", to)
	}
	root.Children(func() {
		kit.Muted(c, o.Rows)
		var names []string
		for _, s := range o.Sizes {
			names = append(names, strconv.Itoa(s))
		}
		pick := strconv.Itoa(*size)
		if choice.Dropdown(c, o.Rows, names, &pick).Changed() {
			*size, _ = strconv.Atoi(pick)
			*page = 0
			root.SetChanged()
			kit.Log("TableFooter", "size", "rows", *size)
		}
		first, last := min(total, *page**size+1), min(total, (*page+1)**size)
		kit.Numeric(c, fmt.Sprintf(o.Range, first, last, total)).TextColor(t.TextMuted)
		for i, step := range []struct {
			name, glyph string
			to          int
		}{{o.Prev, "chevron-left", *page - 1}, {o.Next, "chevron-right", *page + 1}} {
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(step.name).Tooltip(step.name).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)).
				Disabled(!valid(step.to, pages))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
				ui.Icon(c, icon.Must(step.glyph)).Size(16, 16)
			})
			if b.Clicked() {
				turn(step.to)
			}
		}
	})
	return root
}
