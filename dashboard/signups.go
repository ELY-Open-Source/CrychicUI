package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SignupDay is one day's signup count.
type SignupDay struct {
	Label string
	Count int
}

// Signups shows signups by day bars.
func Signups(c *ui.Context, days []SignupDay) ui.Element {
	if len(days) == 0 {
		panic("dashboard: Signups needs days")
	}
	t := c.Theme()
	maxCount := 0
	for _, day := range days {
		if day.Label == "" || day.Count < 0 {
			panic("dashboard: Signups has an invalid day")
		}
		maxCount = max(maxCount, day.Count)
	}
	if maxCount == 0 {
		maxCount = 1
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, day := range days {
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, day.Label).Grow(1)
					kit.Numeric(c, fmt.Sprint(day.Count)).TextColor(t.TextMuted)
				})
				progress(c, float64(day.Count)/float64(maxCount), t.Accent)
			})
		}
	})
	return root
}
