package chart

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// YearHeatOptions tune a YearHeat.
type YearHeatOptions struct {
	// Cell is a day's side; zero gives the default.
	Cell float32
	// Format prints a count; "%g" when empty.
	Format string
	// Date is a Go layout; "Jan 2" default.
	Date string
}

// YearHeat shades a cell per day from start.
func YearHeat(c *ui.Context, label string, start time.Time, counts []float64, at *int, o YearHeatOptions) ui.Element {
	if label == "" || len(counts) == 0 || len(counts) > 371 {
		panic("chart: YearHeat needs a label and 1 to 371 days")
	}
	finite("YearHeat", label, counts...)
	top := 0.0
	for _, v := range counts {
		if v < 0 {
			panic(fmt.Sprintf("chart: YearHeat %q counts must not be negative", label))
		}
		top = max(top, v)
	}
	chosen("YearHeat", label, at, len(counts))
	def(&o.Format, "%g")
	def(&o.Date, "Jan 2")
	t := c.Theme()
	o.Cell = side("YearHeat", label, o.Cell, t.Space(3))
	gap := t.Space(0.5)
	step := o.Cell + gap
	lead := int(start.Weekday())
	cols := (lead + len(counts) + 6) / 7
	w, h := float32(cols)*step-gap, 7*step-gap
	root := ui.Column(c).Gap(t.Space(2))
	return root.Children(func() {
		shown := springs(root, "v", counts)
		wheel{label: label, comp: "YearHeat", n: len(counts), size: w, high: h,
			hit: func(x, y, g float32) int {
				col, row := int((x+w/2)/step), int((y+h/2)/step)
				i := col*7 + row - lead
				if x+w/2 < 0 || y+h/2 < 0 || row > 6 || i < 0 || i >= len(counts) || float32(col) >= g*float32(cols) {
					return -1
				}
				return i
			},
			describe: func(i int) string {
				return start.AddDate(0, 0, i).Format(o.Date) + ": " + fmt.Sprintf(o.Format, counts[i])
			},
			draw: func(p *ui.Painter, r ui.Rect, hot int, g float32) {
				for i, v := range shown {
					k := i + lead
					if float32(k/7) >= g*float32(cols) {
						continue
					}
					ink := t.Border.Alpha(0.5)
					if top > 0 && v > 0 {
						ink = t.Accent.Alpha(0.2 + 0.8*float32(v/top))
					}
					if *at == i || i == hot {
						ink = ink.Mix(t.Text, 0.3)
					}
					p.Fill(ui.Rect{X: r.X + float32(k/7)*step, Y: r.Y + float32(k%7)*step, W: o.Cell, H: o.Cell}, ink, o.Cell/4)
				}
			},
		}.build(c, root, at, nil)
	})
}
