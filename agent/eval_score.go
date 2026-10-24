package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EvalCheck is one scored check of an EvalScore.
type EvalCheck struct {
	Name        string
	Score, Pass float64
}

// EvalScoreOptions word the passed count.
type EvalScoreOptions struct {
	Passed string
}

// EvalScore shows scores per check with pass marks.
func EvalScore(c *ui.Context, label string, checks []EvalCheck, o EvalScoreOptions) ui.Element {
	if label == "" || len(checks) == 0 {
		panic(fmt.Sprintf("agent: EvalScore %q needs a label and checks", label))
	}
	passed := 0
	for _, k := range checks {
		if k.Name == "" || !(k.Score >= 0 && k.Score <= 1) || !(k.Pass >= 0 && k.Pass <= 1) {
			panic(fmt.Sprintf("agent: EvalScore %q check %q needs values from 0 to 1", label, k.Name))
		}
		if k.Score >= k.Pass {
			passed++
		}
	}
	def(&o.Passed, "%d of %d passed")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleList).Label(label)
	g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return root.Children(func() {
		all := passed == len(checks)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			ink := t.Warning
			if all {
				ink = t.Success
			}
			kit.Numeric(c, fmt.Sprintf(o.Passed, passed, len(checks))).FontSize(t.Rem(0.8)).Bold().TextColor(ink).
				Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(ink.Alpha(0.12))
		})
		for i, k := range checks {
			c.Key(k.Name)
			ok := k.Score >= k.Pass
			ink, glyph := t.Danger, "circle-x"
			if ok {
				ink, glyph = t.Success, "circle-check"
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleListItem).Label(k.Name)
			kit.Arrive(row, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*40*time.Millisecond)
			fill := row.AnimateTo("fill", float32(k.Score), ui.Smooth) * g
			row.Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(ink)
				ui.Text(c, k.Name).FontSize(t.Rem(0.85)).Width(t.Space(34)).Shrink(0).SingleLine()
				track := ui.Box(c).Grow(1).Height(t.Space(2)).Radius(999).Background(t.Surface)
				track.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					if fill > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*fill), H: r.H}, ink, 999)
					}
					x := r.X + r.W*float32(k.Pass)
					p.Fill(ui.Rect{X: x - 1, Y: r.Y, W: 2, H: r.H}, t.Text, 0)
				}))
				kit.Numeric(c, fmt.Sprintf("%.2f", k.Score)).FontSize(t.Rem(0.8)).Bold().TextColor(ink).Width(t.Space(10)).Shrink(0)
			})
		}
	})
}
