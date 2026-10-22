package data

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Record is a value at a row and column.
type Record struct {
	Row, Col string
	Value    float64
}

// PivotOptions tune a Pivot.
type PivotOptions struct {
	// Total heads the sums; Format prints a value.
	Total, Format string
}

// Pivot sums records by row and column.
func Pivot(c *ui.Context, label string, records []Record, o PivotOptions) ui.Element {
	if label == "" || len(records) == 0 {
		panic(fmt.Sprintf("data: Pivot %q needs records", label))
	}
	def(&o.Total, "Total")
	def(&o.Format, "%.0f")
	var rows, cols []string
	sums := map[[2]string]float64{}
	for _, r := range records {
		if r.Row == "" || r.Col == "" || math.IsNaN(r.Value) || math.IsInf(r.Value, 0) {
			panic(fmt.Sprintf("data: Pivot %q records need a row, a column and a finite value", label))
		}
		if !slices.Contains(rows, r.Row) {
			rows = append(rows, r.Row)
		}
		if !slices.Contains(cols, r.Col) {
			cols = append(cols, r.Col)
		}
		sums[[2]string{r.Row, r.Col}] += r.Value
	}
	t := c.Theme()
	heads := []Column{{Title: label, Width: t.Space(32)}}
	for _, col := range cols {
		heads = append(heads, Column{Title: col, End: true})
	}
	heads = append(heads, Column{Title: o.Total, End: true})
	// total sums row ("" all) over col (-1 all).
	total := func(row string, col int) float64 {
		s := 0.0
		for _, r := range rows {
			for j, name := range cols {
				if (row == "" || r == row) && (col < 0 || j == col) {
					s += sums[[2]string{r, name}]
				}
			}
		}
		return s
	}
	root := frame(c, label)
	root.Children(func() {
		head(c, root, "Pivot", heads, nil, nil)
		for r, name := range append(slices.Clone(rows), "") {
			c.Key(r)
			line(c, hovered(c, false), heads, nil, func(col int) {
				text := ""
				switch {
				case col == 0 && name == "":
					text = o.Total
				case col == 0:
					text = name
				case col == len(heads)-1:
					text = fmt.Sprintf(o.Format, total(name, -1))
				default:
					text = fmt.Sprintf(o.Format, total(name, col-1))
				}
				cell := ui.Text(c, text).SingleLine()
				if name == "" || col == len(heads)-1 {
					cell.Bold()
				}
				kit.Arrive(cell, ui.Pose{Fade: 1}, time.Duration(r+col)*25*time.Millisecond)
			})
		}
	})
	return root
}
