package health

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Ring is one activity goal; past Goal it laps.
type Ring struct {
	Label, Unit string
	Value, Goal int
}

// ActivityRingsOptions word the legend; Of takes value, goal and unit.
type ActivityRingsOptions struct{ Of string }

// ActivityRings closes three nested rings toward their goals.
func ActivityRings(c *ui.Context, rings [3]Ring, o ActivityRingsOptions) ui.Element {
	labels := make([]string, 3)
	for i, r := range rings {
		if strings.TrimSpace(r.Label) == "" || r.Goal < 1 || r.Value < 0 {
			panic(fmt.Sprintf("health: ActivityRings ring %q needs a label, a goal of 1 or more and a value of 0 or more", r.Label))
		}
		labels[i] = r.Label
	}
	unique("ActivityRings", labels)
	def(&o.Of, "%s / %s %s")
	t := c.Theme()
	inks := [3]ui.Color{t.Danger, t.Success, t.Accent}
	root := card(c).Row().Gap(t.Space(5)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		size := t.Space(32)
		rings3 := ui.Box(c).Size(size, size).Shrink(0)
		var shares [3]float32
		for i, r := range rings {
			grow := rings3.Keyframes(i, []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
				ui.KeyframeOptions{Delay: time.Duration(i) * 150 * time.Millisecond})
			// At most two laps show.
			shares[i] = grow * rings3.AnimateTo([2]int{i, 1}, float32(min(2, float64(r.Value)/float64(r.Goal))), ui.Smooth)
		}
		rings3.Draw(kit.Inside(0, func(p *ui.Painter, b ui.Rect) {
			w := size * 0.11
			cx, cy := b.X+b.W/2, b.Y+b.H/2
			for i := range 3 {
				rr := b.W/2 - w/2 - float32(i)*(w+t.Space(0.75))
				p.StrokePath(ring(cx, cy, rr), w, inks[i].Alpha(0.18))
				first := min(1, shares[i])
				if first > 0 {
					p.StrokePath(ring(cx, cy, rr).Trim(0, first), w, inks[i])
				}
				if lap := shares[i] - 1; lap > 0 {
					p.StrokePath(ring(cx, cy, rr).Trim(0, lap), w, inks[i].Mix(t.Text, 0.25))
				}
			}
		}))
		legend := ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(3))
		legend.Children(func() {
			for i, r := range rings {
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
				kit.Stagger(c, legend, row, i).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(inks[i])
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						kit.Muted(c, r.Label).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprintf(o.Of, number(c, r.Value), number(c, r.Goal), r.Unit)).FontSize(t.Rem(0.9)).Bold().
							TextColor(inks[i]).FillWidth().MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{Numeric: true})
					})
				})
			}
		})
	})
	return root
}
