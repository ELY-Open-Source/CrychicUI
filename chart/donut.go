package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DonutOptions tune a Donut.
type DonutOptions struct {
	Size float32
	// Format prints values; Total names the sum.
	Format, Total string
}

// Donut is a ring with the total in its center.
func Donut(c *ui.Context, label string, names []string, vs []float64, at *int, o DonutOptions) ui.Element {
	def(&o.Format, "%g")
	def(&o.Total, "Total")
	sum := 0.0
	for _, v := range vs {
		sum += v
	}
	t := c.Theme()
	return pies(c, label, names, vs, at, o.Size, 0.62, o.Format, "Donut", func() {
		ui.Column(c).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Format, sum)).FontSize(t.Rem(1.4)).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
			kit.Muted(c, o.Total).FontSize(t.Rem(0.8))
		})
	})
}
