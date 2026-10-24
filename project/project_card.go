package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProjectCardOptions are the project card labels.
type ProjectCardOptions struct{ Members string }

// ProjectCard shows name, members, and progress.
func ProjectCard(c *ui.Context, name string, members []string, value float64, o ProjectCardOptions) ui.Element {
	if name == "" || len(members) == 0 || value < 0 || value > 1 || !finite(value) {
		panic("project: ProjectCard has invalid data")
	}
	def(&o.Members, "%d members")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, name).Bold()
		kit.Muted(c, fmt.Sprintf(o.Members, len(members))).FontSize(t.Rem(0.75))
		progress(c, value, t.Success)
	})
	return root
}
