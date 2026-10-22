package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MatrixOptions tune a Matrix.
type MatrixOptions struct {
	// Mark names a check; it takes the row and column.
	Mark string
}

// Matrix marks row and column pairs with checks.
func Matrix(c *ui.Context, label string, rows, cols []string, marks *ui.Selection[[2]int], o MatrixOptions) ui.Element {
	if label == "" || len(rows) == 0 || len(cols) == 0 || marks == nil {
		panic(fmt.Sprintf("data: Matrix %q needs rows, columns and marks", label))
	}
	unique("Matrix", label, rows)
	unique("Matrix", label, cols)
	def(&o.Mark, "%s, %s")
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
			hovered(c, false).FillWidth().Padding(t.Space(2), t.Space(3)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, name).Width(w).Shrink(0).SingleLine()
				for i, col := range cols {
					c.Key(i)
					ui.Row(c).Grow(1).Basis(0).Justify(ui.Center).Children(func() {
						on := marks.Has([2]int{r, i})
						if checkbox(c, fmt.Sprintf(o.Mark, name, col), &on, false).Changed() {
							toggle(marks, [2]int{r, i})
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("Matrix", "marked", "row", r, "column", i, "on", on)
						}
					})
				}
			})
		}
	})
}
