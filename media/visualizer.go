package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VisualizerOptions tune a Visualizer.
type VisualizerOptions struct {
	// Bars counts the bands; zero gives 48.
	Bars int
	// Size is the ring's width; zero gives the default.
	Size float32
}

// Visualizer rings moving bars around a cover while playing.
func Visualizer(c *ui.Context, label string, art *ui.Bitmap, playing bool, o VisualizerOptions) ui.Element {
	if label == "" || o.Bars < 0 {
		panic(fmt.Sprintf("media: Visualizer %q needs a label and bars not negative", label))
	}
	if o.Bars == 0 {
		o.Bars = 48
	}
	t := c.Theme()
	size("Visualizer", label, &o.Size, t.Space(56))
	box := ui.Box(c).Size(o.Size, o.Size).Center().Role(ui.RoleImage).Label(label)
	hs := bands(c, box, o.Bars, playing)
	r0, reach := o.Size*0.3, o.Size*0.17
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		cx, cy := r.X+r.W/2, r.Y+r.H/2
		bw := float32(2*math.Pi) * r0 / float32(len(hs)) * 0.55
		for i, h := range hs {
			a := 2*math.Pi*float64(i)/float64(len(hs)) - math.Pi/2
			dx, dy := float32(math.Cos(a)), float32(math.Sin(a))
			l := max(t.Space(0.5), h*reach)
			var bar ui.Path
			bar.MoveTo(cx+dx*r0, cy+dy*r0).LineTo(cx+dx*(r0+l), cy+dy*(r0+l))
			p.StrokePath(&bar, bw, t.Accent.Mix(t.Success, float32(i)/float32(len(hs))))
		}
	}))
	return box.Children(func() {
		cover(c, art, r0*1.76).Radius(999).AnimatePose(ui.Pose{Scale: 1 + 0.04*unit(playing)}, ui.Snappy)
	})
}
