package feedback

import "github.com/egoist/mygo/ui"

// SpinSize is a Spin's size step.
type SpinSize uint8

// The three Spin sizes.
const (
	SpinMedium SpinSize = iota
	SpinSmall
	SpinLarge
)

// SpinOptions tune a Spin.
type SpinOptions struct {
	// Size picks small, medium or large; medium when zero.
	Size SpinSize
}

// Spin is an arc turning around a faint ring.
func Spin(c *ui.Context, label string, o SpinOptions) ui.Element {
	if label == "" {
		panic("feedback: Spin needs a label")
	}
	t := c.Theme()
	var size float32
	switch o.Size {
	case SpinSmall:
		size = t.Space(4)
	case SpinMedium:
		size = t.Space(6)
	case SpinLarge:
		size = t.Space(10)
	default:
		panic("feedback: unknown SpinSize")
	}
	w := max(1.5, size/10)
	return ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleProgress).Label(label).Range(0, 1, -1).Draw(func(p *ui.Painter, r ui.Rect) {
		p.AnimationFrame()
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
		turn := float32(p.Now().UnixMilli()%800) / 800 * 360
		var ring, arc ui.Path
		p.StrokePath(ring.Circle(cx, cy, rad), w, t.Accent.Alpha(0.18))
		p.Transform(ui.Rotation(turn), cx, cy, func() { p.StrokePath(arc.Circle(cx, cy, rad).Trim(0, 0.28), w, t.Accent) })
	})
}
