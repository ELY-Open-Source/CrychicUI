package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AccentCard is a card with an accent border growing in.
func AccentCard(c *ui.Context, content func()) ui.Element {
	t := c.Theme()
	root := kit.Card(c).FillWidth().Clip()
	return root.Children(func() {
		bar := ui.Box(c).Absolute().Left(0).Top(0).Bottom(0).Width(t.Space(1)).Background(t.Accent).TransformOrigin(0.5, 0)
		p := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		bar.Scale(1, p)
		ui.Column(c).Padding(t.Space(4), t.Space(4), t.Space(4), t.Space(4)+t.Space(3)).Children(content)
	})
}
