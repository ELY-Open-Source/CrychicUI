package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Plan shows steps with pending, running and done marks.
func Plan(c *ui.Context, label string, steps []string, current int) ui.Element {
	if label == "" || len(steps) == 0 || current < 0 || current >= len(steps) {
		panic(fmt.Sprintf("agent: Plan %q needs a label, steps and a current step", label))
	}
	for _, s := range steps {
		if s == "" {
			panic(fmt.Sprintf("agent: Plan %q steps must not be empty", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, s := range steps {
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			row.Children(func() {
				ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Children(func() {
					switch {
					case i < current:
						done := ui.Icon(c, icon.Must("circle-check")).Size(16, 16).TextColor(t.Success)
						kit.Arrive(done, ui.Pose{Scale: 0.5, Fade: 1}, 0)
					case i == current:
						ic := ui.Icon(c, icon.Must("rotate-cw")).Size(16, 16).TextColor(t.Accent)
						ic.Rotate(360 * float32(ic.Loop("spin", time.Second, ui.Linear)))
					default:
						ui.Icon(c, icon.Must("circle")).Size(16, 16).TextColor(t.TextMuted)
					}
				})
				text := ui.Text(c, s).SingleLine()
				switch {
				case i < current:
					text.TextColor(t.TextMuted)
				case i == current:
					text.Bold()
				default:
					text.TextColor(t.TextMuted)
				}
			})
		}
	})
}
