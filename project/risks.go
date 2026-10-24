package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Risk is one risk by likelihood and impact.
type Risk struct {
	Title              string
	Likelihood, Impact int
}

// Risks shows likelihood by impact grid.
func Risks(c *ui.Context, risks []Risk) ui.Element {
	if len(risks) == 0 {
		panic("project: Risks needs risks")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(5).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for impact := 1; impact <= 5; impact++ {
			for likelihood := 1; likelihood <= 5; likelihood++ {
				c.Key(fmt.Sprintf("%d.%d", impact, likelihood))
				found := false
				cell := ui.Column(c)
				for _, risk := range risks {
					if risk.Title == "" || risk.Likelihood < 1 || risk.Likelihood > 5 || risk.Impact < 1 || risk.Impact > 5 {
						panic("project: Risks has an invalid risk")
					}
					if risk.Likelihood != likelihood || risk.Impact != impact {
						continue
					}
					found = true
					cell.Children(func() {
						ui.Text(c, risk.Title).Center().Padding(t.Space(1)).Radius(t.Space(0.5)).
							Background(t.Danger.Alpha(float32(likelihood*impact) / 25)).FontSize(t.Rem(0.65))
					})
				}
				if !found {
					cell.Children(func() { ui.Text(c, "") })
				}
			}
		}
	})
	return root
}
