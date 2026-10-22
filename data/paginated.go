package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PaginatedOptions tune a Paginated.
type PaginatedOptions struct {
	// Size is rows per page; zero gives 10.
	Size int
	// Of takes the page and the page count.
	Prev, Next, Of string
}

// Paginated pages a table from 0; pages cross-fade.
func Paginated(c *ui.Context, label string, cols []Column, n int, cell func(row, col int), page *int, o PaginatedOptions) ui.Element {
	columns("Paginated", label, cols)
	if n < 0 || cell == nil || page == nil || o.Size < 0 {
		panic(fmt.Sprintf("data: Paginated %q needs a row count, cells, a page and a size", label))
	}
	if o.Size == 0 {
		o.Size = 10
	}
	pages := max(1, (n+o.Size-1)/o.Size)
	if !valid(*page, pages) {
		panic(fmt.Sprintf("data: Paginated %q page %d is not among %d", label, *page, pages))
	}
	def(&o.Prev, "Previous page")
	def(&o.Next, "Next page")
	def(&o.Of, "Page %d of %d")
	t := c.Theme()
	root := frame(c, label)
	shown := *page
	fade := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{}, Exit: &ui.Motion{}}
	return root.Children(func() {
		head(c, root, "Paginated", cols, nil, nil)
		c.Key(shown)
		ui.Column(c).FillWidth().Transition(fade).Children(func() {
			for i := shown * o.Size; i < min(n, (shown+1)*o.Size); i++ {
				c.Key(i)
				line(c, hovered(c, false), cols, nil, func(col int) { cell(i, col) })
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.End).Padding(t.Space(2), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Of, shown+1, pages)).TextColor(t.TextMuted)
			for i, step := range []struct {
				name, glyph string
				to          int
			}{{o.Prev, "chevron-left", shown - 1}, {o.Next, "chevron-right", shown + 1}} {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(step.name).Tooltip(step.name).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius)).
					Disabled(!valid(step.to, pages))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Icon(c, icon.Must(step.glyph)).Size(16, 16)
				})
				if b.Clicked() {
					*page = step.to
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("Paginated", "page", "index", step.to)
				}
			}
		})
	})
}
