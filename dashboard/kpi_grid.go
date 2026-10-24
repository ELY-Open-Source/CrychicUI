package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KPI is one tile; Icon is Lucide, Points a trend.
type KPI struct {
	Label, Value, Delta, Icon string
	Points                    []float64
}

// KPIGrid shows KPI tiles in two columns; values roll.
func KPIGrid(c *ui.Context, tiles []KPI) ui.Element {
	if len(tiles) == 0 {
		panic("dashboard: KPIGrid needs tiles")
	}
	names := make([]string, len(tiles))
	for i, tile := range tiles {
		if tile.Value == "" || len(tile.Points) == 1 {
			panic(fmt.Sprintf("dashboard: KPIGrid tile %q needs a value and no single point", tile.Label))
		}
		names[i] = tile.Label
	}
	unique("KPIGrid", names)
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, tile := range tiles {
			c.Key(tile.Label)
			box := kit.Card(c).Gap(t.Space(1.5)).Padding(t.Space(3)).MinWidth(0)
			box.BorderColor(box.AnimateColor("edge", t.Border.Mix(t.Accent, 0.5*unit(box.Hovered())), ui.Snappy))
			kit.Stagger(c, root, box, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					if tile.Icon != "" {
						ui.Icon(c, icon.Must(tile.Icon)).Size(14, 14).Shrink(0)
					}
					kit.Muted(c, tile.Label).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
				})
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, tile.Value).Bold().FontSize(t.Rem(1.5)).Grow(1).MinWidth(0).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
					if tile.Delta != "" {
						deltaChip(c, tile.Delta)
					}
				})
				if len(tile.Points) > 1 {
					points := make([]Point, len(tile.Points))
					for j, value := range tile.Points {
						points[j] = Point{Value: value}
					}
					series(c, fmt.Sprintf("KPIGrid tile %q", tile.Label), points, t.Accent).Height(t.Space(10))
				}
			})
		}
	})
	return root
}
