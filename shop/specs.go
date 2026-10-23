package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SpecRow is one specification name and value.
type SpecRow struct {
	Name, Value string
}

// SpecGroup is one specification group.
type SpecGroup struct {
	Name string
	Rows []SpecRow
}

// Specs shows specification rows by group.
func Specs(c *ui.Context, groups []SpecGroup) ui.Element {
	if len(groups) == 0 {
		panic("shop: Specs needs groups")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, group := range groups {
			if group.Name == "" || len(group.Rows) == 0 {
				panic("shop: Specs group needs a name and rows")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Text(c, group.Name).Bold()
				for j, row := range group.Rows {
					if row.Name == "" || row.Value == "" {
						panic("shop: Specs row must not be empty")
					}
					c.Key(j)
					ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
						kit.Muted(c, row.Name).FontSize(t.Rem(0.8)).Grow(1)
						ui.Text(c, row.Value).FontSize(t.Rem(0.85)).Wrap()
					})
				}
			})
		}
	})
	return root
}
