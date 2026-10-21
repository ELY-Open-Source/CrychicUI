package typography

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// Blockquote is a quote beside a growing accent bar.
func Blockquote(c *ui.Context, quote string) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Children(func() {
		bar := ui.Box(c).Width(t.Space(0.5)).FillHeight().Background(t.Accent).TransformOrigin(0.5, 0)
		grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		bar.Scale(1, grow)
		ui.Text(c, quote)
	})
}
