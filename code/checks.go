package code

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ChecksOptions word the tally and a running check.
type ChecksOptions struct {
	Passing, Running string
}

// Checks lists CI checks with spinners, results and run times.
func Checks(c *ui.Context, label string, checks []Check, o ChecksOptions) ui.Element {
	if label == "" || len(checks) == 0 {
		panic(fmt.Sprintf("code: Checks %q needs a label and checks", label))
	}
	seen := map[string]bool{}
	for _, k := range checks {
		checkOne("Checks", k)
		if seen[k.Name] {
			panic(fmt.Sprintf("code: Checks %q needs unique names", label))
		}
		seen[k.Name] = true
	}
	def(&o.Passing, "%d of %d passing")
	def(&o.Running, "Running…")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	root.Children(func() {
		panelHead(c, "shield-check", label, func() { checkTally(c, checks, o.Passing) })
		ui.Column(c).FillWidth().Padding(t.Space(1)).Children(func() {
			for i, k := range checks {
				c.Key(k.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
					Role(ui.RoleListItem).Label(k.Name)
				row.Background(row.AnimateColor("bg", t.Danger.Alpha(0.06*unit(k.State == CheckFail)), ui.Snappy))
				stagger(c, root, row, i).Children(func() {
					switch k.State {
					case CheckPending:
						e := ui.Icon(c, iconSpin).Size(16, 16).TextColor(t.Warning)
						e.Rotate(e.Loop("spin", 1200*time.Millisecond, ui.Linear) * 360)
					case CheckPass:
						ui.Icon(c, iconCheckCircle).Size(16, 16).TextColor(t.Success).ContentTransition(ui.ContentTransition{})
					case CheckFail:
						ui.Icon(c, iconXCircle).Size(16, 16).TextColor(t.Danger).ContentTransition(ui.ContentTransition{})
					}
					ui.Text(c, k.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					when := o.Running
					if k.State != CheckPending {
						when = took(k.Duration)
					}
					kit.Numeric(c, when).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0).ContentTransition(ui.ContentTransition{})
				})
			}
		})
	})
	return root
}
