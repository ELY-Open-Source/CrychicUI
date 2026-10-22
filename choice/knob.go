package choice

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// KnobOptions tune a Knob.
type KnobOptions struct {
	// Step is one detent; a hundredth of the range.
	Step float64
	// Format writes the value; "%.0f" when empty.
	Format string
}

// Knob is a dial dragged or keyed; its arc fills.
func Knob(c *ui.Context, label string, value *float64, lo, hi float64, o KnobOptions) ui.Element {
	sliderArgs("Knob", label, value, lo, hi, o.Step)
	if o.Step == 0 {
		o.Step = (hi - lo) / 100
	}
	if o.Format == "" {
		o.Format = "%.0f"
	}
	t := c.Theme()
	size := t.Space(16)
	root := ui.Box(c).Size(size, size).Center().Focusable().Label(label).Role(ui.RoleSlider).Cursor(ui.CursorResizeNS)
	drag := ui.Local(root, "drag", func() float64 { return *value })
	set := func(v float64) {
		v = snap(min(hi, max(lo, v)), lo, o.Step)
		if v > hi {
			v -= o.Step
		}
		if v != *value {
			*value = v
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Knob", "changed", "value", v)
		}
	}
	if _, dy, held := root.Dragged(); held {
		if dy != 0 {
			*drag -= float64(dy) / float64(2*size) * (hi - lo)
			*drag = min(hi, max(lo, *drag))
			set(*drag)
		}
	} else {
		*drag = *value
	}
	up, down := ui.KeyRight, ui.KeyLeft
	if c.Direction() == ui.RightToLeft {
		up, down = down, up
	}
	switch {
	case root.Shortcut(0, up), root.Shortcut(0, ui.KeyUp):
		set(*value + o.Step)
	case root.Shortcut(0, down), root.Shortcut(0, ui.KeyDown):
		set(*value - o.Step)
	case root.Shortcut(0, ui.KeyHome):
		set(lo)
	case root.Shortcut(0, ui.KeyEnd):
		set(hi)
	}
	root.Range(lo, hi, *value).Step(o.Step)
	frac := root.AnimateTo("frac", fraction(*value, lo, hi), ui.Snappy)
	w := t.Space(1.5)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w
		p.Transform(ui.Rotation(135), cx, cy, func() {
			var track, fill ui.Path
			p.StrokePath(track.Circle(cx, cy, rad).Trim(0, 0.75), w, t.Border)
			if frac > 0 {
				p.StrokePath(fill.Circle(cx, cy, rad).Trim(0, 0.75*frac), w, t.Accent)
			}
		})
		a := (135 + 270*float64(frac)) * math.Pi / 180
		x, y, d := cx+rad*float32(math.Cos(a)), cy+rad*float32(math.Sin(a)), w*2.4
		p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.AccentText, d/2)
		p.Stroke(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Accent, d/2, 1)
	})
	return root.Children(func() {
		kit.Numeric(c, fmt.Sprintf(o.Format, *value)).Bold()
	})
}
