package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FunnelStage is one conversion stage.
type FunnelStage struct {
	Label string
	Count int
}

// Conversion shows a funnel with rates.
func Conversion(c *ui.Context, stages []FunnelStage) ui.Element {
	if len(stages) < 2 || stages[0].Count <= 0 {
		panic("dashboard: Conversion needs stages")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	previous := stages[0].Count
	root.Children(func() {
		for i, stage := range stages {
			if stage.Label == "" || stage.Count < 0 || stage.Count > previous {
				panic("dashboard: Conversion has invalid stages")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				rate := 0.0
				if previous > 0 {
					rate = float64(stage.Count) / float64(previous)
				}
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, stage.Label).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%d (%.0f%%)", stage.Count, 100*rate)).TextColor(t.TextMuted)
				})
				progress(c, float64(stage.Count)/float64(stages[0].Count), t.Accent)
			})
			previous = stage.Count
		}
	})
	return root
}
