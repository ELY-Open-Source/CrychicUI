package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FocusBlock is one focus interval as day fractions.
type FocusBlock struct {
	Start, End float64
}

// FocusTime shows focus blocks in a day bar.
func FocusTime(c *ui.Context, blocks []FocusBlock) ui.Element {
	if len(blocks) == 0 {
		panic("schedule: FocusTime needs blocks")
	}
	t := c.Theme()
	for _, block := range blocks {
		if block.Start < 0 || block.End <= block.Start || block.End > 1 || !finite(block.Start) || !finite(block.End) {
			panic("schedule: FocusTime has an invalid block")
		}
	}
	root := ui.Box(c).FillWidth().Height(t.Space(8)).Radius(t.Radius).Clip().Background(t.Surface)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for _, block := range blocks {
			x := r.X + r.W*float32(block.Start)
			w := r.W * float32(block.End-block.Start)
			p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, t.Accent.Alpha(0.25), t.Radius)
		}
	}))
	return root
}
