package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SpecRow is one specification name and value.
type SpecRow struct {
	Name, Value string
}

// SpecGroup is one specification group; Icon is a Lucide name.
type SpecGroup struct {
	Name, Icon string
	Rows       []SpecRow
}

// Specs shows specification rows under group headings.
func Specs(c *ui.Context, groups []SpecGroup) ui.Element {
	if len(groups) == 0 {
		panic("shop: Specs needs groups")
	}
	for _, g := range groups {
		if g.Name == "" || len(g.Rows) == 0 {
			panic(fmt.Sprintf("shop: Specs group %q needs a name and rows", g.Name))
		}
		for _, r := range g.Rows {
			if r.Name == "" || r.Value == "" {
				panic(fmt.Sprintf("shop: Specs group %q row %q needs a name and value", g.Name, r.Name))
			}
		}
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, g := range groups {
			c.Key(i)
			section := ui.Column(c).FillWidth().Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(groups)-1)))
			kit.Stagger(c, root, section, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, 0, t.Space(1.5), 0).Children(func() {
					if g.Icon != "" {
						ui.Icon(c, icon.Must(g.Icon)).Size(16, 16).TextColor(t.Accent)
					}
					ui.Text(c, g.Name).Bold().Shrink(0)
				})
				for j, r := range g.Rows {
					c.Key(j)
					ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(1.5), 0).AlignItems(ui.Start).
						BorderWidth(1, 0, 0, 0).BorderColor(t.Border.Alpha(0.6)).Children(func() {
						kit.Muted(c, r.Name).FontSize(t.Rem(0.8)).Width(t.Space(32)).Shrink(0).Wrap()
						ui.Text(c, r.Value).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
					})
				}
			})
		}
	})
	return root
}
