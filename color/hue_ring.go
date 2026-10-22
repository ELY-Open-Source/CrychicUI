package color

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HueRing sets a hue by dragging round a ring.
func HueRing(c *ui.Context, label string, hue *float64) ui.Element {
	if label == "" || hue == nil || !(*hue >= 0 && *hue < 360) {
		panic(fmt.Sprintf("color: HueRing %q needs a hue from 0 to 360", label))
	}
	t := c.Theme()
	size, w := t.Space(36), t.Space(4)
	root := ui.Box(c).Size(size, size).Shrink(0).Focusable().Role(ui.RoleSlider).Label(label).Cursor(ui.CursorPointer)
	step := ui.Local(root, "step", func() int { return int(*hue / 15) })
	set := func(h float64, keyed bool) {
		if h = math.Mod(h+360, 360); h == *hue {
			return
		}
		*hue = h
		root.SetChanged()
		if s := int(h / 15); keyed || s != *step {
			*step = s
			sound.Play(sound.Tick)
			kit.Log("HueRing", "turned", "step", s)
		}
	}
	switch {
	case root.Shortcut(0, ui.KeyRight) || root.Shortcut(0, ui.KeyUp):
		set(math.Floor(*hue)+1, true)
	case root.Shortcut(0, ui.KeyLeft) || root.Shortcut(0, ui.KeyDown):
		set(math.Ceil(*hue)-1, true)
	}
	_, _, held := root.Dragged()
	if x, y, _ := root.PointerPosition(); held {
		set(math.Atan2(float64(x-size/2), float64(size/2-y))*180/math.Pi, false)
	}
	root.Range(0, 360, *hue).Step(1)
	knob := (w/2 - 1) * (0.87 + 0.13*root.AnimateTo("grip", unit(held), ui.Snappy))
	h := *hue
	stops := make([]ui.GradientStop, 7)
	for i := range stops {
		stops[i] = ui.GradientStop{Color: fromHSV(float64(i)*60, 1, 1, 255), At: float32(i) / 6}
	}
	return kit.Arrive(root, ui.Pose{Scale: 0.9, Fade: 1}, 0).Draw(kit.Inside(size/2, func(p *ui.Painter, r ui.Rect) {
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
		var band ui.Path
		p.StrokePathConicGradient(band.Circle(cx, cy, rad), w, ui.ConicGradient{Stops: stops})
		var core ui.Path
		p.FillPath(core.Circle(cx, cy, rad-w), fromHSV(h, 1, 1, 255))
		a := h * math.Pi / 180
		kx, ky := cx+rad*float32(math.Sin(a)), cy-rad*float32(math.Cos(a))
		var k ui.Path
		p.FillPath(k.Circle(kx, ky, knob), t.Background)
		p.StrokePath(&k, 2, t.Text)
	}))
}
