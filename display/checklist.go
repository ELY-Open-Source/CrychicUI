package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChecklistItem is one read-only checked row.
type ChecklistItem struct {
	Label string
	Done  bool
}

// Checklist fades rows in and grows their checks.
func Checklist(c *ui.Context, items []ChecklistItem) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
		for i, item := range items {
			row := ui.Row(c).Key(i).Gap(t.Space(2)).AlignItems(ui.Center)
			kit.Arrive(row, ui.Pose{Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				mark := ui.Column(c).Size(t.Space(5), t.Space(5)).Center().Radius(t.Space(1)).Border(1, t.Border)
				mark.Children(func() {
					check := ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success)
					p := check.AnimateTo("draw", unit(item.Done), ui.Snappy)
					check.Scale(p, p).Opacity(p)
				})
				ui.Text(c, item.Label)
			})
		}
	})
}
