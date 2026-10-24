package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RunStep is one step of a RunTimeline.
type RunStep struct {
	Name     string
	Duration time.Duration
	State    StepState
}

// RunTimelineOptions word the summary.
type RunTimelineOptions struct {
	Format string
}

// RunTimeline shows steps of a run as growing waterfall bars.
func RunTimeline(c *ui.Context, label string, steps []RunStep, o RunTimelineOptions) ui.Element {
	if label == "" || len(steps) == 0 {
		panic(fmt.Sprintf("agent: RunTimeline %q needs a label and steps", label))
	}
	total := time.Duration(0)
	for _, s := range steps {
		if s.Name == "" || s.Duration < 0 || s.State > StepFailed {
			panic(fmt.Sprintf("agent: RunTimeline %q step %q needs a name, a duration and a known state", label, s.Name))
		}
		total += s.Duration
	}
	if total <= 0 {
		panic(fmt.Sprintf("agent: RunTimeline %q needs a positive total", label))
	}
	def(&o.Format, "%d steps · %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleList).Label(label)
	g := root.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	return root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, label).Bold().Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Format, len(steps), clock(total))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
		})
		start := time.Duration(0)
		for i, s := range steps {
			c.Key(i)
			from, span := float32(start)/float32(total), float32(s.Duration)/float32(total)
			start += s.Duration
			ink := t.Accent
			switch s.State {
			case StepDone:
				ink = t.Success
			case StepFailed:
				ink = t.Danger
			case StepPending:
				ink = t.TextMuted
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleListItem).Label(s.Name)
			kit.Arrive(row, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(i)*50*time.Millisecond)
			x := row.AnimateTo("from", from, ui.Smooth)
			w := row.AnimateTo("span", span, ui.Smooth) * g
			ink = row.AnimateColor("ink", ink, ui.Snappy)
			row.Children(func() {
				mark(c, s.State)
				ui.Text(c, s.Name).FontSize(t.Rem(0.85)).Width(t.Space(26)).Shrink(0).SingleLine()
				track := ui.Box(c).Grow(1).Height(t.Space(2)).Radius(999).Background(t.Surface)
				track.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
					if w > 0 {
						p.Fill(ui.Rect{X: r.X + r.W*x, Y: r.Y, W: max(r.H, r.W*w), H: r.H}, ink, 999)
					}
				}))
				kit.Numeric(c, clock(s.Duration)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).Width(t.Space(9)).Shrink(0)
			})
		}
	})
}
