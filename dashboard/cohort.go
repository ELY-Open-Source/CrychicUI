package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CohortRow is one cohort, its size and retention by period.
type CohortRow struct {
	Label  string
	Size   int
	Values []float64
}

// CohortOptions word the grid; Period and Users take numbers.
type CohortOptions struct{ Title, Period, Users string }

// Cohort shades a retention grid; hover reads a cell.
func Cohort(c *ui.Context, rows []CohortRow, o CohortOptions) ui.Element {
	if len(rows) == 0 {
		panic("dashboard: Cohort needs rows")
	}
	names, cols := make([]string, len(rows)), 0
	for i, row := range rows {
		if row.Size < 0 || len(row.Values) == 0 || len(row.Values) > 12 {
			panic(fmt.Sprintf("dashboard: Cohort %q needs a size of zero or more and 1–12 values", row.Label))
		}
		for _, v := range row.Values {
			if v < 0 || v > 1 || !finite(v) {
				panic(fmt.Sprintf("dashboard: Cohort %q values must be in [0,1]", row.Label))
			}
		}
		names[i] = row.Label
		cols = max(cols, len(row.Values))
	}
	unique("Cohort", names)
	def(&o.Title, "Retention")
	def(&o.Period, "Week %d")
	def(&o.Users, "%s users")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "grid-3x3", o.Title, nil)
		body := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3))
		body.Children(func() {
			hover, at := -1, -1
			cell := t.Space(11)
			ui.ScrollHorizontal(c).FillWidth().Children(func() {
				grid := ui.Column(c).FillWidth().MinWidth(t.Space(22) + float32(cols)*(cell+t.Space(0.5))).Gap(t.Space(0.5))
				grid.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
						ui.Box(c).Width(t.Space(22)).Shrink(0)
						for j := range cols {
							c.Key(j)
							kit.Muted(c, fmt.Sprintf(o.Period, j)).FontSize(t.Rem(0.65)).Grow(1).Basis(0).MinWidth(cell).TextAlign(ui.Center).SingleLine()
						}
					})
					for i, row := range rows {
						c.Key(row.Label)
						line := ui.Row(c).FillWidth().Gap(t.Space(0.5)).AlignItems(ui.Stretch)
						kit.Stagger(c, grid, line, i).Children(func() {
							ui.Column(c).Width(t.Space(22)).Shrink(0).Justify(ui.Center).Children(func() {
								ui.Text(c, row.Label).Bold().FontSize(t.Rem(0.75)).FillWidth().SingleLine()
								kit.Muted(c, fmt.Sprintf(o.Users, number(c, row.Size))).FontSize(t.Rem(0.65)).FillWidth().SingleLine()
							})
							for j := range cols {
								c.Key(j)
								if j >= len(row.Values) {
									ui.Box(c).Grow(1).Basis(0).MinWidth(cell)
									continue
								}
								v := row.Values[j]
								box := ui.Box(c).Grow(1).Basis(0).MinWidth(cell).Center().Padding(t.Space(1.5), 0).Radius(t.Space(0.75)).
									Border(1, ui.Transparent)
								if box.Hovered() {
									hover, at = i, j
								}
								box.Background(box.AnimateColor("shade", t.Accent.Alpha(0.08+0.82*float32(v)), ui.Snappy)).
									BorderColor(box.AnimateColor("edge", t.Text.Alpha(0.6*unit(box.Hovered())), ui.Snappy))
								ink := t.Text
								if v > 0.55 {
									ink = t.AccentText
								}
								box.Children(func() {
									kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*v)).FontSize(t.Rem(0.7)).TextColor(ink).Shrink(0)
								})
							}
						})
					}
				})
			})
			ui.Row(c).FillWidth().Children(func() {
				if hover < 0 {
					kit.Muted(c, " ").FontSize(t.Rem(0.75))
					return
				}
				row := rows[hover]
				text := fmt.Sprintf("%s · %s · %.0f%%", row.Label, fmt.Sprintf(o.Period, at), 100*row.Values[at])
				kit.Muted(c, text).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap().ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root
}
