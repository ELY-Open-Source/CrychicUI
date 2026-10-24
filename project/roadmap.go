package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Quarter is one roadmap quarter and initiatives.
type Quarter struct {
	Label string
	Items []string
}

// Roadmap shows quarters with initiatives.
func Roadmap(c *ui.Context, quarters []Quarter) ui.Element {
	if len(quarters) == 0 {
		panic("project: Roadmap needs quarters")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, quarter := range quarters {
			if quarter.Label == "" || len(quarter.Items) == 0 {
				panic("project: Roadmap has an invalid quarter")
			}
			c.Key(i)
			ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Children(func() {
				ui.Text(c, quarter.Label).Bold()
				for _, item := range quarter.Items {
					if item == "" {
						panic("project: Roadmap initiative must not be empty")
					}
					ui.Text(c, item).FontSize(t.Rem(0.8)).Wrap()
				}
			})
		}
	})
	return root
}
