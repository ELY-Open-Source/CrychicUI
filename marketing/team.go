package marketing

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Member is one person; nil Photo shows initials.
type Member struct {
	Name, Role string
	Photo      *ui.Bitmap
}

// TeamOptions word the grid; zero Columns is three.
type TeamOptions struct {
	Title, Text string
	Columns     int
}

// Team lays people out in a grid of lifting cards.
func Team(c *ui.Context, members []Member, o TeamOptions) ui.Element {
	if o.Columns == 0 {
		o.Columns = 3
	}
	if len(members) == 0 || o.Columns < 1 || o.Columns > 4 {
		panic(fmt.Sprintf("marketing: Team %q needs members and 1 to 4 columns", o.Title))
	}
	names := make([]string, len(members))
	for i, m := range members {
		names[i] = m.Name
	}
	unique("Team", names)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		grid := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(3))
		grid.Children(func() {
			for i, m := range members {
				c.Key(strconv.Quote(m.Name))
				cell := ui.Column(c).MinWidth(0)
				kit.Stagger(c, grid, cell, i).Children(func() {
					face := ui.Column(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4), t.Space(2)).
						Radius(t.Radius * 1.5)
					face.Border(1, face.AnimateColor("edge", t.Border.Mix(t.Accent, 0.5*unit(face.Hovered())), ui.Snappy))
					kit.Lift(face).Children(func() {
						avatar(c, m.Name, m.Photo, 16)
						ui.Text(c, m.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).Wrap()
						if m.Role != "" {
							kit.Muted(c, m.Role).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).TextAlign(ui.Center).Wrap()
						}
					})
				})
			}
		})
	})
	return root
}
