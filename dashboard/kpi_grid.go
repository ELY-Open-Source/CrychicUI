package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KPI is one metric tile.
type KPI struct {
	Label, Value, Delta string
	Points              []float64
}

// KPIGrid shows KPI tiles with rolling values.
func KPIGrid(c *ui.Context, tiles []KPI) ui.Element {
	if len(tiles) == 0 {
		panic("dashboard: KPIGrid needs tiles")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, tile := range tiles {
			if tile.Label == "" || tile.Value == "" || tile.Delta == "" {
				panic("dashboard: KPIGrid has an invalid tile")
			}
			c.Key(i)
			card := ui.Column(c).Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
			card.Children(func() {
				kit.Muted(c, tile.Label).FontSize(t.Rem(0.75))
				kit.Numeric(c, tile.Value).Bold().FontSize(t.Rem(1.3))
				ui.Text(c, tile.Delta).FontSize(t.Rem(0.75)).TextColor(t.Accent)
				if len(tile.Points) > 0 {
					points := make([]Point, len(tile.Points))
					for j, value := range tile.Points {
						points[j] = Point{Value: value}
					}
					series(c, points, t.Accent)
				}
			})
		}
	})
	return root
}
