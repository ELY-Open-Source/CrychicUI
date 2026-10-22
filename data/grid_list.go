package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Tile is a label and Lucide icon.
type Tile struct{ Label, Icon string }

// GridListOptions tune a GridList.
type GridListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
}

// GridList is a ui.GridView of chosen tiles.
func GridList(c *ui.Context, label string, tiles []Tile, chosen *int, o GridListOptions) ui.Element {
	if label == "" || len(tiles) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(tiles) {
		panic(fmt.Sprintf("data: GridList %q needs tiles and a choice or -1", label))
	}
	for _, tl := range tiles {
		if tl.Label == "" {
			panic(fmt.Sprintf("data: GridList %q tiles need labels", label))
		}
		icon.Must(tl.Icon)
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Height(height(c, o.Height))
	grid := ui.Local(root, "grid", func() ui.GridState { return ui.GridState{} })
	grid.Selected = chosen
	grid.Label = func(i int) string { return tiles[i].Label }
	root.Children(func() {
		g := ui.GridView(c, grid, len(tiles), t.Space(28), t.Space(24), func(i int) {
			ui.Column(c).FillWidth().FillHeight().Center().Gap(t.Space(1.5)).Children(func() {
				badge := ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(999).Background(t.Accent.Alpha(0.12))
				k := 1 + 0.12*badge.AnimateTo("pop", unit(i == *chosen), ui.Bouncy)
				badge.Scale(k, k).Children(func() { ui.Icon(c, icon.Must(tiles[i].Icon)).Size(16, 16).TextColor(t.Accent) })
				ui.Text(c, tiles[i].Label).FontSize(t.Rem(0.8)).SingleLine()
			})
		}).Label(label).Gap(t.Space(1)).Grow(1)
		if g.Changed() {
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("GridList", "chose", "index", *chosen)
		}
	})
	return root
}
