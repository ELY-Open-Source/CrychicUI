// Code components for the gallery.
package code

import (
	"github.com/egoist/mygo/ui"
	"time"
)

// Checks are CI checks with spinners and results.
func Checks(c *ui.Context, label string, checks []Check) ui.Element {
	if label == "" || len(checks) == 0 {
		panic("code: Checks needs a label and checks")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(0.5)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for _, k := range checks {
			checkOne("Checks", k)
			c.Key(k.Name)
			row := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(0.5), 0)
			row.Children(func() {
				switch k.State {
				case CheckPending:
					e := ui.Icon(c, iconSpin).Size(16, 16).TextColor(t.Accent)
					e.Rotate(e.Loop("spin", 1200*time.Millisecond, ui.Linear) * 360)
				case CheckPass:
					ui.Icon(c, iconCheckCircle).Size(16, 16).TextColor(t.Success)
				case CheckFail:
					ui.Icon(c, iconXCircle).Size(16, 16).TextColor(t.Danger)
				}
				ui.Text(c, k.Name).SingleLine()
			})
		}
	})
	return root
}

// GraphCommit is one commit in a graph.
