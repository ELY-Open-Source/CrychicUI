package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RegionValue is one region and value; Code tags it.
type RegionValue struct {
	Name, Code string
	Value      float64
}

// RegionsOptions word the card; Total takes the sum.
type RegionsOptions struct{ Title, Total string }

// Regions ranks regions with value bars and shares.
func Regions(c *ui.Context, regions []RegionValue, o RegionsOptions) ui.Element {
	if len(regions) == 0 {
		panic("dashboard: Regions needs regions")
	}
	names, peak, total := make([]string, len(regions)), 0.0, 0.0
	for i, region := range regions {
		if region.Value < 0 || !finite(region.Value) {
			panic(fmt.Sprintf("dashboard: Regions %q needs a finite value of zero or more", region.Name))
		}
		names[i] = region.Name
		peak = max(peak, region.Value)
		total += region.Value
	}
	unique("Regions", names)
	if !finite(total) {
		panic("dashboard: Regions total overflows")
	}
	def(&o.Title, "Regions")
	def(&o.Total, "%s total")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "map-pin", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Total, number(c, total)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			for i, region := range regions {
				c.Key(region.Name)
				line := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, line, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						if region.Code != "" {
							ui.Box(c).Shrink(0).Padding(t.Space(0.25), t.Space(1.5)).Radius(t.Space(1)).Background(t.Surface).Children(func() {
								ui.Text(c, region.Code).Font("monospace").FontSize(t.Rem(0.7)).Bold().Shrink(0)
							})
						}
						ui.Text(c, region.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, number(c, region.Value)).Bold().FontSize(t.Rem(0.85)).Shrink(0)
						share := 0.0
						if total > 0 {
							share = region.Value / total
						}
						kit.Muted(c, fmt.Sprintf("%.0f%%", 100*share)).FontSize(t.Rem(0.75)).Width(t.Space(9)).TextAlign(ui.End).Shrink(0)
					})
					fill := 0.0
					if peak > 0 {
						fill = region.Value / peak
					}
					progress(c, fill, t.Accent)
				})
			}
		})
	})
	return root
}
