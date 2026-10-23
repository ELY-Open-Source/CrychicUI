package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Leader is one performer and score.
type Leader struct {
	Name  string
	Score int
}

// Leaders shows top performers with ranks.
func Leaders(c *ui.Context, people []Leader) ui.Element {
	if len(people) == 0 {
		panic("dashboard: Leaders needs people")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, person := range people {
			if person.Name == "" || person.Score < 0 {
				panic("dashboard: Leaders has an invalid person")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprintf("#%d", i+1)).TextColor(t.Accent).Width(t.Space(6))
				ui.Text(c, person.Name).Grow(1)
				kit.Numeric(c, fmt.Sprint(person.Score)).Bold()
			})
		}
	})
	return root
}
