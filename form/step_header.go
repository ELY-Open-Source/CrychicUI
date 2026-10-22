package form

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/egoist/mygo/ui"
)

// StepHeader numbers steps; connectors fill and the current step pops.
func StepHeader(c *ui.Context, steps []string, current int) ui.Element {
	if len(steps) == 0 || current < 0 || current >= len(steps) {
		panic(fmt.Sprintf("form: StepHeader step %d of %d", current, len(steps)))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
	last := ui.Local(root, "current", func() int { return current })
	run := ui.Local(root, "run", func() uint64 { return 0 })
	if current != *last {
		*last = current
		(*run)++
	}
	return root.Children(func() {
		for i, step := range steps {
			if i > 0 {
				line := ui.Box(c).Grow(1).Height(t.Space(0.5)).Radius(999)
				fill := line.AnimateTo("fill", unit(i <= current), ui.Smooth)
				line.Draw(func(p *ui.Painter, r ui.Rect) {
					p.Fill(r, t.Border, r.H/2)
					p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, t.Accent, r.H/2)
				})
			}
			reached := i <= current
			dot := ui.Box(c).Size(t.Space(7), t.Space(7)).Shrink(0).Radius(999).Center().Border(1, t.Border)
			bg, ink := t.Background, t.TextMuted
			if reached {
				bg, ink = t.Accent, t.AccentText
			}
			dot.Background(dot.AnimateColor("bg", bg, ui.Snappy)).TextColor(ink)
			if i == current && *run > 0 {
				k := dot.Keyframes("pop", []ui.Keyframe{{Value: 1}, {At: 120 * time.Millisecond, Value: 1.2},
					{At: 280 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}, ui.KeyframeOptions{Run: *run})
				dot.Scale(k, k)
			}
			dot.Children(func() {
				if i < current {
					ui.Icon(c, icon.Must("check")).Size(16, 16)
				} else {
					ui.Text(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.85)).FontWeight(600)
				}
			})
			name := ui.Text(c, step).SingleLine()
			if !reached {
				name.TextColor(t.TextMuted)
			}
		}
	})
}
