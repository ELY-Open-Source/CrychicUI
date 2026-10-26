package marketing

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Cell is a plan's entry or its mark.
type Cell struct {
	Included bool
	Text     string
}

// CompareRow is one feature across the plans.
type CompareRow struct {
	Feature string
	Cells   []Cell
}

// PlanCompareOptions word the matrix; Highlight tints one plan.
type PlanCompareOptions struct{ Title, Highlight, Yes, No string }

// PlanCompare is a feature matrix; it scrolls when narrow.
func PlanCompare(c *ui.Context, plans []string, rows []CompareRow, o PlanCompareOptions) ui.Element {
	def(&o.Title, "Compare plans")
	def(&o.Yes, "Included")
	def(&o.No, "Not included")
	if len(plans) == 0 || len(plans) > 6 || len(rows) == 0 || o.Highlight != "" && !slices.Contains(plans, o.Highlight) {
		panic(fmt.Sprintf("marketing: PlanCompare %q needs 1 to 6 plans, rows and a highlight among the plans", o.Title))
	}
	unique("PlanCompare", plans)
	features := make([]string, len(rows))
	for i, r := range rows {
		if len(r.Cells) != len(plans) {
			panic(fmt.Sprintf("marketing: PlanCompare row %q has %d cells for %d plans", r.Feature, len(r.Cells), len(plans)))
		}
		features[i] = r.Feature
	}
	unique("PlanCompare", features)
	t := c.Theme()
	cell, name, gap := t.Space(20), t.Space(36), t.Space(1)
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "table-2", o.Title, nil)
		view := ui.ScrollHorizontal(c).FillWidth()
		seen := view.LayoutBounds().W
		kit.Settle(c, view, []float32{seen}, func() []float32 { return []float32{view.LayoutBounds().W} })
		view.Children(func() {
			grid := ui.Column(c).Width(max(seen, name+float32(len(plans))*(cell+gap)+t.Space(6))).Padding(t.Space(1.5), t.Space(3))
			grid.Children(func() {
				ui.Row(c).FillWidth().Gap(gap).Padding(t.Space(1.5), 0).Children(func() {
					ui.Box(c).Width(name).Shrink(0)
					for i, p := range plans {
						c.Key(i)
						lit := p == o.Highlight
						ui.Text(c, p).FontSize(t.Rem(0.8)).Bold().TextColor(t.Text.Mix(t.Accent, unit(lit))).Grow(1).Basis(0).MinWidth(cell).
							TextAlign(ui.Center).SingleLine()
					}
				})
				for i, r := range rows {
					c.Key(strconv.Quote(r.Feature))
					line := ui.Row(c).FillWidth().Gap(gap).AlignItems(ui.Center).Padding(t.Space(2), 0).Radius(t.Radius)
					line.Background(line.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.5*unit(i%2 == 0)+0.5*unit(line.Hovered())), ui.Snappy))
					kit.Stagger(c, grid, line, i).Children(func() {
						ui.Text(c, r.Feature).FontSize(t.Rem(0.8)).Width(name).Shrink(0).Padding(0, t.Space(2)).Wrap()
						for j, x := range r.Cells {
							c.Key(j)
							box := ui.Row(c).Grow(1).Basis(0).MinWidth(cell).Justify(ui.Center).AlignItems(ui.Center)
							box.Children(func() {
								switch {
								case x.Text != "":
									ui.Text(c, x.Text).FontSize(t.Rem(0.78)).Bold().MinWidth(0).SingleLine()
								case x.Included:
									glyph(c, "check", 16).TextColor(t.Success).Label(o.Yes)
								default:
									glyph(c, "minus", 16).TextColor(t.TextMuted.Alpha(0.6)).Label(o.No)
								}
							})
						}
					})
				}
			})
		})
	})
	return root
}
