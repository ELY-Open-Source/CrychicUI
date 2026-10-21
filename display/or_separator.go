package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OrSeparator puts a label between two growing lines.
func OrSeparator(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		line := func(key string) {
			bar := ui.Box(c).Key(key).Grow(1).Basis(0).Height(1).Background(t.Border)
			p := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			bar.Scale(p, 1)
		}
		line("left")
		kit.Muted(c, label).FontSize(t.Rem(0.85))
		line("right")
	})
}
