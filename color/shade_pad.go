package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShadePad sets saturation across and value up.
func ShadePad(c *ui.Context, label string, value *ui.Color) ui.Element {
	if label == "" || value == nil {
		panic("color: ShadePad needs a label and a color")
	}
	t := c.Theme()
	w, h := t.Space(48), t.Space(32)
	root := ui.Box(c).Size(w, h).Shrink(0).Radius(t.Radius).Focusable().Role(ui.RoleSlider).Label(label).Cursor(ui.CursorPointer)
	vals := ui.Local(root, "hsv", func() [3]float64 { h, s, v := hsv(*value); return [3]float64{h, s, v} })
	made := ui.Local(root, "made", func() ui.Color { return *value })
	if *value != *made {
		h, s, v := hsv(*value)
		if s == 0 || v == 0 {
			h = vals[0]
		}
		*vals, *made = [3]float64{h, s, v}, *value
	}
	bucket := ui.Local(root, "bucket", func() [2]int { return [2]int{int(vals[1] * 10), int(vals[2] * 10)} })
	set := func(ns, nv float64, keyed bool) {
		ns, nv = max(0, min(ns, 1)), max(0, min(nv, 1))
		if ns == vals[1] && nv == vals[2] {
			return
		}
		vals[1], vals[2] = ns, nv
		*value = fromHSV(vals[0], ns, nv, value.A)
		*made = *value
		root.SetChanged()
		if b := [2]int{int(ns * 10), int(nv * 10)}; keyed || b != *bucket {
			*bucket = b
			sound.Play(sound.Tick)
			kit.Log("ShadePad", "moved", "s", b[0], "v", b[1])
		}
	}
	switch {
	case root.Shortcut(0, ui.KeyRight):
		set(vals[1]+0.01, vals[2], true)
	case root.Shortcut(0, ui.KeyLeft):
		set(vals[1]-0.01, vals[2], true)
	case root.Shortcut(0, ui.KeyUp):
		set(vals[1], vals[2]+0.01, true)
	case root.Shortcut(0, ui.KeyDown):
		set(vals[1], vals[2]-0.01, true)
	}
	_, _, held := root.Dragged()
	if x, y, _ := root.PointerPosition(); held {
		set(float64(x/w), float64(1-y/h), false)
	}
	s, v := vals[1], vals[2]
	root.Range(0, 100, s*100)
	grip := root.AnimateTo("grip", unit(held), ui.Snappy)
	full, white, black := fromHSV(vals[0], 1, 1, 255), fromHSV(0, 0, 1, 255), fromHSV(0, 0, 0, 255)
	return kit.Arrive(root, ui.Pose{Scale: 0.96, Fade: 1}, 0).Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		p.FillGradient(r, ui.LinearGradient{From: white, To: full, Angle: 90}, 0)
		p.FillGradient(r, ui.LinearGradient{From: black.Alpha(0), To: black, Angle: 180}, 0)
		var k ui.Path
		p.StrokePath(k.Circle(r.X+r.W*float32(s), r.Y+r.H*float32(1-v), t.Space(2)*(1+0.2*grip)), 2, white)
	}))
}
