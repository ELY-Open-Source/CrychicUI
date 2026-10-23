package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RegionValue is one region and value.
type RegionValue struct {
	Name  string
	Value float64
}

// Regions shows regions with value bars.
func Regions(c *ui.Context, regions []RegionValue) ui.Element {
	if len(regions) == 0 {
		panic("dashboard: Regions needs regions")
	}
	t := c.Theme()
	maxValue := 0.0
	for _, region := range regions {
		if region.Name == "" || region.Value < 0 || !finite(region.Value) {
			panic("dashboard: Regions has an invalid region")
		}
		maxValue = max(maxValue, region.Value)
	}
	if maxValue == 0 {
		maxValue = 1
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, region := range regions {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, region.Name).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f", region.Value)).TextColor(t.TextMuted)
				})
				progress(c, region.Value/maxValue, t.Accent)
			})
		}
	})
	return root
}
