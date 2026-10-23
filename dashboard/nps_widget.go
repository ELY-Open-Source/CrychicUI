package dashboard

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NPSWidgetOptions are the NPS labels.
type NPSWidgetOptions struct{ Score, Promoters, Passives, Detractors string }

// NPSWidget shows promoters, passives, and detractors.
func NPSWidget(c *ui.Context, promoters, passives, detractors int, o NPSWidgetOptions) ui.Element {
	if promoters < 0 || passives < 0 || detractors < 0 {
		panic("dashboard: NPSWidget has invalid counts")
	}
	def(&o.Score, "NPS %d")
	def(&o.Promoters, "Promoters")
	def(&o.Passives, "Passives")
	def(&o.Detractors, "Detractors")
	t := c.Theme()
	total := float64(promoters) + float64(passives) + float64(detractors)
	if total == 0 {
		panic("dashboard: NPSWidget has invalid counts")
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Score, int(math.Round(100*(float64(promoters)-float64(detractors))/total)))).Bold().FontSize(t.Rem(1.4))
		for _, row := range []struct {
			name  string
			count int
			ink   ui.Color
		}{{o.Promoters, promoters, t.Success}, {o.Passives, passives, t.Warning}, {o.Detractors, detractors, t.Danger}} {
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, row.name).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%d (%.0f%%)", row.count, 100*float64(row.count)/total)).TextColor(t.TextMuted)
				})
				progress(c, float64(row.count)/total, row.ink)
			})
		}
	})
	return root
}
