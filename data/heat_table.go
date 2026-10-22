package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HeatTableOptions tune a HeatTable.
type HeatTableOptions struct {
	// Format prints a value; "%g" when empty.
	Format string
}

// HeatTable tints cells by value; tints animate.
func HeatTable(c *ui.Context, label string, rows, cols []string, values [][]float64, o HeatTableOptions) ui.Element {
	if label == "" || len(rows) == 0 || len(cols) == 0 || len(values) != len(rows) {
		panic(fmt.Sprintf("data: HeatTable %q needs rows, columns and a value row each", label))
	}
	unique("HeatTable", label, rows)
	unique("HeatTable", label, cols)
	lo, hi := values[0][0], values[0][0]
	for _, vs := range values {
		if len(vs) != len(cols) {
			panic(fmt.Sprintf("data: HeatTable %q rows need a value per column", label))
		}
		for _, v := range vs {
			if !finite64(v) {
				panic(fmt.Sprintf("data: HeatTable %q values must be finite", label))
			}
			lo, hi = min(lo, v), max(hi, v)
		}
	}
	def(&o.Format, "%g")
	t := c.Theme()
	w := t.Space(28)
	root := frame(c, label)
	return root.Children(func() {
		ui.Row(c).FillWidth().Background(t.Surface).Padding(t.Space(2), t.Space(3)).Children(func() {
			ui.Box(c).Width(w).Shrink(0)
			for i, col := range cols {
				c.Key(i)
				kit.Muted(c, col).Bold().FontSize(t.Rem(0.8)).Grow(1).Basis(0).TextAlign(ui.Center).SingleLine()
			}
		})
		for r, name := range rows {
			c.Key(name)
			ui.Row(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(0.5), t.Space(3)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, name).Width(w).Shrink(0).SingleLine()
				for i, v := range values[r] {
					c.Key(i)
					k := float32(0.5)
					if hi > lo {
						k = float32((v - lo) / (hi - lo))
					}
					cell := ui.Box(c).Grow(1).Basis(0).Center().Padding(t.Space(2), t.Space(1)).Radius(t.Radius / 2)
					cell.Background(cell.AnimateColor("bg", t.Accent.Alpha(0.06+0.7*k), ui.Smooth)).
						TextColor(cell.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(k > 0.5)), ui.Smooth)).
						Label(fmt.Sprintf("%s, %s", name, cols[i])).Children(func() {
						kit.Numeric(c, fmt.Sprintf(o.Format, v))
					})
				}
			})
		}
	})
}
