package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

// TotalsRowOptions tune a TotalsRow.
type TotalsRowOptions struct {
	// Total heads the first column.
	Total string
}

// TotalsRow is a footer of column sums; digits roll.
func TotalsRow(c *ui.Context, label string, cols []Column, sums []int64, o TotalsRowOptions) ui.Element {
	columns("TotalsRow", label, cols)
	if len(sums) != len(cols)-1 {
		panic(fmt.Sprintf("data: TotalsRow %q needs a sum per column after the first", label))
	}
	def(&o.Total, "Total")
	t := c.Theme()
	root := ui.Row(c).Background(t.Surface).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Role(ui.RoleGroup).Label(label)
	return line(c, root, cols, nil, func(col int) {
		if col == 0 {
			ui.Text(c, o.Total).Bold()
			return
		}
		motion.RollingDigits(c, sums[col-1]).FontSize(t.Rem(1))
	})
}
