package feedback

import (
	"time"

	"github.com/egoist/mygo/ui"
)

var grow = []ui.Keyframe{{Value: 0}, {At: 320 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}

// Callout is a tinted box whose left border grows in.
func Callout(c *ui.Context, tone Tone, title string, content func()) ui.Element {
	if content == nil {
		panic("feedback: Callout needs content")
	}
	t := c.Theme()
	col, _ := look(t, tone)
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3), t.Space(4)).PaddingStart(t.Space(5)).
		Radius(t.Radius).Background(col.Alpha(0.06))
	return root.Children(func() {
		edge := ui.Box(c).Absolute().Start(0).Top(0).Bottom(0).Width(t.Space(1)).Radius(t.Radius, 0, 0, t.Radius).
			Background(col).TransformOrigin(0.5, 0)
		edge.Scale(1, edge.Keyframes("grow", grow))
		if title != "" {
			ui.Text(c, title).Bold().TextColor(t.Text.Mix(col, 0.3))
		}
		content()
	})
}
