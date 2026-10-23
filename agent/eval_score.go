package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EvalCheck is one scored check of an EvalScore.
type EvalCheck struct {
	Name        string
	Score, Pass float64
}

// EvalScore shows scores per check with pass marks.
func EvalScore(c *ui.Context, label string, checks []EvalCheck) ui.Element {
	if label == "" || len(checks) == 0 {
		panic(fmt.Sprintf("agent: EvalScore %q needs a label and checks", label))
	}
	for _, k := range checks {
		if k.Name == "" || !(k.Score >= 0 && k.Score <= 1) || !(k.Pass >= 0 && k.Pass <= 1) {
			panic(fmt.Sprintf("agent: EvalScore %q checks need a name and values from 0 to 1", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup).Label(label)
	g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return root.Children(func() {
		for i, k := range checks {
			c.Key(i)
			ink := t.Danger
			if k.Score >= k.Pass {
				ink = t.Success
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			fill := row.AnimateTo(fmt.Sprint("fill.", i), float32(k.Score), ui.Smooth) * g
			row.Children(func() {
				ui.Text(c, k.Name).Width(t.Space(24)).Shrink(0).SingleLine()
				track := ui.Box(c).Grow(1).Height(t.Space(2)).Radius(999).Background(t.Surface).Clip()
				track.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					if fill > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, ink, 999)
					}
					x := r.X + r.W*float32(k.Pass)
					p.Line(x, r.Y-t.Space(0.5), x, r.Y+r.H+t.Space(0.5), t.Space(0.5), t.Text)
				}))
				kit.Numeric(c, fmt.Sprintf("%.2f", k.Score)).FontSize(t.Rem(0.8)).TextColor(ink).Width(t.Space(10)).Shrink(0)
			})
		}
	})
}
