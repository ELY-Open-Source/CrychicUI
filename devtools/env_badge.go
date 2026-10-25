package devtools

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EnvBadgeOptions name the four stages.
type EnvBadgeOptions struct{ Stages [4]string }

// EnvBadge is an environment pill colored by stage.
func EnvBadge(c *ui.Context, stage Stage, o EnvBadgeOptions) ui.Element {
	validStage("EnvBadge", stage)
	for i, name := range [...]string{"Development", "Preview", "Staging", "Production"} {
		def(&o.Stages[i], name)
	}
	t := c.Theme()
	ink := stageInk(t, stage)
	root := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.75), t.Space(2.5)).Radius(999).Shrink(0).
		Role(ui.RoleStatus)
	root.Background(root.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).Border(1, root.AnimateColor("line", ink.Alpha(0.4), ui.Snappy))
	kit.Arrive(root, ui.Pose{Scale: 0.85, Fade: 1}, 0)
	root.Children(func() {
		c.Key(stage)
		dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
		kit.Arrive(dot, ui.Pose{Scale: 2.4, Fade: 1}, 120*time.Millisecond)
		ui.Text(c, o.Stages[stage]).FontSize(t.Rem(0.75)).Bold().TextColor(ink).Shrink(0).ContentTransition(ui.ContentTransition{})
	})
	return root
}
