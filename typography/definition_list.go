package typography

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Definition pairs a term with its description.
type Definition struct{ Term, Description string }

// DefinitionList fades term and description rows in.
func DefinitionList(c *ui.Context, items []Definition) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(3)).Children(func() {
		for i, item := range items {
			row := ui.Row(c).Key(i).Gap(t.Space(3))
			kit.Arrive(row, ui.Pose{Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Text(c, item.Term).Bold().Width(t.Rem(10))
				ui.Text(c, item.Description).Grow(1).ContentTransition(ui.ContentTransition{})
			})
		}
	})
}
