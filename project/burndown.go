package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BurndownPoint is one remaining and ideal pair.
type BurndownPoint struct {
	Day       string
	Remaining float64
	Ideal     float64
}

// Burndown shows remaining against the ideal line.
func Burndown(c *ui.Context, points []BurndownPoint) ui.Element {
	if len(points) < 2 {
		panic("project: Burndown needs points")
	}
	t := c.Theme()
	maxValue := 0.0
	for _, point := range points {
		if point.Day == "" || point.Remaining < 0 || point.Ideal < 0 || !finite(point.Remaining) || !finite(point.Ideal) {
			panic("project: Burndown has an invalid point")
		}
		maxValue = max(maxValue, point.Remaining, point.Ideal)
	}
	if maxValue == 0 {
		maxValue = 1
	}
	root := ui.Box(c).FillWidth().Height(t.Space(24))
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		var actual, ideal ui.Path
		for i, point := range points {
			x := r.X + r.W*float32(i)/float32(len(points)-1)
			ya := r.Y + r.H*(1-float32(point.Remaining/maxValue))
			yi := r.Y + r.H*(1-float32(point.Ideal/maxValue))
			if i == 0 {
				actual.MoveTo(x, ya)
				ideal.MoveTo(x, yi)
			} else {
				actual.LineTo(x, ya)
				ideal.LineTo(x, yi)
			}
		}
		p.StrokePath(&actual, 2, t.Accent)
		p.StrokePath(&ideal, 1, t.TextMuted.Alpha(0.5))
	}))
	return root
}
