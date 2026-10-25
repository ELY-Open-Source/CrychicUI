package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QueryResultOptions word the grid; Row and Rows take a count.
type QueryResultOptions struct {
	Title, Row, Rows, Empty, Null string
	Height                        float32
}

// QueryResult is a result grid with count and time.
func QueryResult(c *ui.Context, columns []string, rows [][]string, elapsed time.Duration, o QueryResultOptions) ui.Element {
	def(&o.Title, "Result")
	if len(columns) == 0 || elapsed < 0 || o.Height < 0 {
		panic(fmt.Sprintf("devtools: QueryResult %q needs columns and a non-negative time", o.Title))
	}
	for i, r := range rows {
		if len(r) != len(columns) {
			panic(fmt.Sprintf("devtools: QueryResult %q row %d has %d cells for %d columns", o.Title, i, len(r), len(columns)))
		}
	}
	def(&o.Row, "%s row")
	def(&o.Rows, "%s rows")
	def(&o.Empty, "No rows")
	def(&o.Null, "NULL")
	t := c.Theme()
	if o.Height == 0 {
		o.Height = t.Space(60)
	}
	cell, gap := t.Space(20), t.Space(3)
	// A column is numeric when every non-null cell parses.
	numeric := make([]bool, len(columns))
	for j := range columns {
		numeric[j] = len(rows) > 0
		for _, r := range rows {
			if _, err := strconv.ParseFloat(r[j], 64); err != nil && r[j] != o.Null {
				numeric[j] = false
			}
		}
	}
	root := card(c).Role(ui.RoleTable)
	root.Children(func() {
		kit.CardHead(c, "table", o.Title, func() {
			format := o.Rows
			if len(rows) == 1 {
				format = o.Row
			}
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				glyph(c, "timer", 12).TextColor(t.TextMuted)
				kit.Muted(c, took(elapsed)).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
			})
			kit.Chip(c, fmt.Sprintf(format, number(c, len(rows))), t.Accent)
		})
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		view.Children(func() {
			w := max(seen, float32(len(columns))*(cell+gap)-gap+t.Space(6))
			ui.Column(c).Width(w).Children(func() {
				line := func(cells []string, head bool) {
					ui.Row(c).FillWidth().Gap(gap).Padding(t.Space(1.5), t.Space(3)).Children(func() {
						for j, v := range cells {
							resultCell(c, v, cell, head, numeric[j], o.Null)
						}
					})
				}
				ui.Column(c).FillWidth().Background(t.SurfaceHover.Alpha(0.6)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).
					Children(func() { line(columns, true) })
				if len(rows) == 0 {
					empty(c, "table", o.Empty)
					return
				}
				body := ui.Scroll(c).FillWidth().MaxHeight(o.Height)
				body.Children(func() {
					list := ui.Column(c).FillWidth()
					list.Children(func() {
						for i, r := range rows {
							c.Key(i)
							stripe := ui.Column(c).FillWidth()
							if i%2 == 1 {
								stripe.Background(t.SurfaceHover.Alpha(0.35))
							}
							kit.Stagger(c, list, stripe, i).Children(func() { line(r, false) })
						}
					})
				})
			})
		})
	})
	return root
}

// resultCell is one grid cell; numeric columns align right.
func resultCell(c *ui.Context, v string, w float32, head, numeric bool, null string) {
	t := c.Theme()
	at := ui.Start
	if numeric {
		at = ui.End
	}
	ui.Row(c).Grow(1).Basis(0).MinWidth(w).Justify(at).Children(func() {
		text := ui.Text(c, v).Font("monospace").FontSize(t.Rem(0.72)).MinWidth(0).SingleLine()
		switch {
		case head:
			text.Bold().TextColor(t.TextMuted)
		case v == null:
			text.TextColor(t.TextMuted.Alpha(0.7))
		case numeric:
			text.TextColor(t.Accent)
		}
	})
}
