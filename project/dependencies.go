package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DependenciesOptions are the dependency labels.
type DependenciesOptions struct{ BlockedBy, Blocks string }

// Dependencies shows blocked-by and blocks links.
func Dependencies(c *ui.Context, blockedBy, blocks []string, opened *string, o DependenciesOptions) ui.Element {
	if opened == nil || len(blockedBy) == 0 && len(blocks) == 0 {
		panic("project: Dependencies needs links and state")
	}
	def(&o.BlockedBy, "Blocked by")
	def(&o.Blocks, "Blocks")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, group := range []struct {
			name  string
			links []string
		}{{o.BlockedBy, blockedBy}, {o.Blocks, blocks}} {
			c.Key(i)
			ui.Column(c).Gap(t.Space(1)).Children(func() {
				kit.Muted(c, group.name).FontSize(t.Rem(0.75))
				for j, link := range group.links {
					if link == "" {
						panic("project: Dependencies link must not be empty")
					}
					c.Key(j)
					b := kit.Press(ui.ButtonBase(c).Label(link).Role(ui.RoleLink).Padding(t.Space(0.5), t.Space(1)).Radius(t.Radius))
					b.Children(func() { ui.Text(c, link).TextColor(t.Accent) })
					if b.Clicked() {
						*opened = link
						root.SetChanged()
						sound.Play(sound.Open)
						kit.Log("Dependencies", "opened")
					}
				}
			})
		}
	})
	return root
}
