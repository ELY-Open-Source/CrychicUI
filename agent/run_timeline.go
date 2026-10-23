package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RunStep is one line of a RunTimeline.
type RunStep struct {
	Name     string
	Duration time.Duration
}

// RunTimeline shows steps of a run with growing duration bars.
func RunTimeline(c *ui.Context, label string, steps []RunStep) ui.Element {
	if label == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: RunTimeline %q needs a label and steps", label))
	}
	longest := time.Duration(0)
	for _, s := range steps {
		if s.Name == "" || s.Duration < 0 {
			panic(fmt.Sprintf("agent: RunTimeline %q steps need a name and a duration", label))
		}
		longest = max(longest, s.Duration)
	}
	if longest == 0 {
		panic(fmt.Sprintf("agent: RunTimeline %q needs a positive duration", label))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup).Label(label)
	g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 500 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return root.Children(func() {
		for i, s := range steps {
			c.Key(i)
			frac := float32(s.Duration) / float32(longest)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			bar := row.AnimateTo(fmt.Sprint("bar.", i), frac, ui.Smooth) * g
			row.Children(func() {
				ui.Text(c, s.Name).Width(t.Space(24)).Shrink(0).SingleLine()
				track := ui.Box(c).Grow(1).Height(t.Space(1.5)).Radius(999).Background(t.Surface).Clip()
				track.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					if bar > 0 {
						p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * bar, H: r.H}, t.Accent, 999)
					}
				}))
				kit.Muted(c, clock(s.Duration)).FontSize(t.Rem(0.8)).Width(t.Space(8)).Shrink(0)
			})
		}
	})
}
