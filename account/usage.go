package account

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Resource is one metered account resource; Icon is optional.
type Resource struct {
	Name, Unit, Icon string
	Used, Limit      float64
}

// UsageOptions word the card; Format takes used, limit and unit.
type UsageOptions struct{ Title, Format, Over string }

// Usage shows usage meters by resource.
func Usage(c *ui.Context, resources []Resource, o UsageOptions) ui.Element {
	if len(resources) == 0 {
		panic("account: Usage needs resources")
	}
	names := make([]string, len(resources))
	for i, r := range resources {
		if r.Unit == "" || r.Used < 0 || r.Limit <= 0 || !finite(r.Used) || !finite(r.Limit) {
			panic(fmt.Sprintf("account: Usage %q needs a unit, a used amount of at least 0 and a positive limit", r.Name))
		}
		names[i] = r.Name
	}
	unique("Usage", names)
	def(&o.Title, "Usage this month")
	def(&o.Format, "%s of %s %s")
	def(&o.Over, "Over limit")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "gauge", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(3.5)).Padding(t.Space(4))
		list.Children(func() {
			for i, r := range resources {
				share := r.Used / r.Limit
				if !finite(share) {
					panic(fmt.Sprintf("account: Usage %q share overflows", r.Name))
				}
				pct := math.Round(share * 100)
				if !finite(pct) {
					panic(fmt.Sprintf("account: Usage %q percent overflows", r.Name))
				}
				if share < 1 {
					pct = min(pct, 99)
				}
				ink := t.Accent
				switch {
				case share >= 1:
					ink = t.Danger
				case share >= 0.8:
					ink = t.Warning
				}
				c.Key(r.Name)
				col := ui.Column(c).FillWidth().Gap(t.Space(1.5))
				kit.Stagger(c, list, col, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						mark := r.Icon
						def(&mark, "gauge")
						glyph(c, mark, 16).TextColor(t.TextMuted)
						ui.Text(c, r.Name).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(t.Space(20)).SingleLine()
						used := strconv.FormatFloat(r.Used, 'g', -1, 64)
						limit := strconv.FormatFloat(r.Limit, 'g', -1, 64)
						kit.Numeric(c, fmt.Sprintf(o.Format, used, limit, r.Unit)).FontSize(t.Rem(0.78)).TextColor(t.TextMuted).MinWidth(0).SingleLine()
					})
					meter(c, float32(min(1, share)), ink)
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						pct := kit.Numeric(c, strconv.FormatFloat(pct, 'f', 0, 64)+"%").FontSize(t.Rem(0.72)).Bold().Grow(1)
						pct.TextColor(pct.AnimateColor("ink", ink, ui.Snappy))
						if share > 1 {
							kit.Chip(c, o.Over, t.Danger)
						}
					})
				})
			}
		})
	})
	return root
}
