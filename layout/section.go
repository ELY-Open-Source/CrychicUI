package layout

import (
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Section is a muted heading, growing divider and body.
func Section(c *ui.Context, title string, content func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).Gap(t.Space(2)).Children(func() {
		kit.Muted(c, title).Bold().FontSize(t.Rem(0.85))
		line := ui.Box(c).FillWidth().Height(1).Background(t.Border).TransformOrigin(0, 0.5)
		grow := line.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		line.Scale(grow, 1)
		ui.Column(c).Gap(t.Space(2)).Children(content)
	})
}
