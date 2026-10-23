package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DayDivider is a day label between two lines; fades in.
func DayDivider(c *ui.Context, label string) ui.Element {
	if label == "" {
		panic("chat: DayDivider needs a label")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		ui.Box(c).Grow(1).Height(1).Background(t.Border)
		kit.Muted(c, label).FontSize(t.Rem(0.8)).Bold().SingleLine()
		ui.Box(c).Grow(1).Height(1).Background(t.Border)
	})
}
