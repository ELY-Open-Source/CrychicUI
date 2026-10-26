package marketing

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Feature is an icon, a title and a line.
type Feature struct{ Icon, Title, Text string }

// FeatureGridOptions word the grid; zero Columns is two.
type FeatureGridOptions struct {
	Title, Text string
	Columns     int
}

// FeatureGrid lays features out in a grid; they stagger in.
func FeatureGrid(c *ui.Context, features []Feature, o FeatureGridOptions) ui.Element {
	if o.Columns == 0 {
		o.Columns = 2
	}
	if len(features) == 0 || o.Columns < 1 || o.Columns > 4 {
		panic(fmt.Sprintf("marketing: FeatureGrid %q needs features and 1 to 4 columns", o.Title))
	}
	titles := make([]string, len(features))
	for i, f := range features {
		if f.Icon == "" {
			panic(fmt.Sprintf("marketing: FeatureGrid feature %q needs an icon", f.Title))
		}
		titles[i] = f.Title
	}
	unique("FeatureGrid", titles)
	t := c.Theme()
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		grid := ui.Grid(c).Columns(o.Columns).FillWidth().Gap(t.Space(3))
		grid.Children(func() {
			for i, f := range features {
				c.Key(strconv.Quote(f.Title))
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius * 1.5)
				cell.Border(1, cell.AnimateColor("edge", t.Border.Mix(inks[i%4], 0.5*unit(cell.Hovered())), ui.Snappy))
				cell.Background(cell.AnimateColor("bg", t.Background.Mix(inks[i%4], 0.04*unit(cell.Hovered())), ui.Snappy))
				kit.Stagger(c, grid, cell, i).Children(func() {
					tile(c, f.Icon, inks[i%4], t.Space(10))
					ui.Text(c, f.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).Wrap()
					kit.Muted(c, f.Text).FontSize(t.Rem(0.8)).FillWidth().MinWidth(0).Wrap()
				})
			}
		})
	})
	return root
}
