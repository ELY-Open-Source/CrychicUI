package choice

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChoiceMatrix is radio rows across shared columns; -1 unanswered.
func ChoiceMatrix(c *ui.Context, rows, columns []string, selected *[]int) ui.Element {
	distinct("ChoiceMatrix", rows)
	distinct("ChoiceMatrix", columns)
	if selected == nil || len(*selected) != len(rows) {
		panic("choice: ChoiceMatrix needs one answer per row")
	}
	for _, v := range *selected {
		if v < -1 || v >= len(columns) {
			panic("choice: ChoiceMatrix answer out of range")
		}
	}
	t := c.Theme()
	cell := t.Space(22)
	root := ui.Column(c).FillWidth().Dividers(1, t.Border).Gap(1)
	return root.Children(func() {
		ui.Row(c).Padding(t.Space(2), 0).Children(func() {
			ui.Spacer(c)
			for _, col := range columns {
				c.Key(col)
				kit.Muted(c, col).FontSize(t.Rem(0.85)).Width(cell).TextAlign(ui.Center)
			}
		})
		for r, name := range rows {
			c.Key(name)
			ui.Row(c).Padding(t.Space(1), 0).Children(func() {
				ui.Text(c, name).Grow(1)
				v := (*selected)[r]
				changed := false
				ui.RadioGroup(c, func() {
					for i, col := range columns {
						c.Key(col)
						b := ui.RadioBase(c, &v, i).Label(col).Width(cell).Height(t.Space(8)).Center().Radius(t.Radius)
						if b.Hovered() {
							b.Background(t.SurfaceHover)
						}
						b.Children(func() { dot(c, b, v == i) })
						changed = kit.Changed(b, sound.Select) || changed
					}
				}).Label(name).Row()
				if changed {
					(*selected)[r] = v
					root.SetChanged()
					kit.Log("ChoiceMatrix", "answered", "row", r, "column", v)
				}
			})
		}
	})
}
