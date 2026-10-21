package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// StepList shows done, current and next steps read-only.
func StepList(c *ui.Context, steps []string, current int) ui.Element {
	if len(steps) == 0 || current < 0 || current >= len(steps) {
		panic("display: StepList needs steps and a valid current index")
	}
	t := c.Theme()
	return ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
		for i, step := range steps {
			ui.Row(c).Key(i).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				mark := ui.Column(c).Size(t.Space(5), t.Space(5)).Center().Radius(999).Border(1, t.Border)
				bg := t.Background
				if i == current {
					bg = t.Accent
				}
				mark.Background(mark.AnimateColor("bg", bg, ui.Snappy))
				mark.Children(func() {
					if i < current {
						ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success).ContentTransition(ui.ContentTransition{})
					} else if i == current {
						ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Accent).ContentTransition(ui.ContentTransition{})
					}
				})
				ui.Text(c, step)
			})
		}
	})
}
