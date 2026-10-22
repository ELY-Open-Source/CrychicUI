package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ZoomImageOptions tune a ZoomImage.
type ZoomImageOptions struct {
	// Max zoom; zero gives 4.
	Max float64
	// Height of the view; zero gives the default.
	Height  float32
	In, Out string
}

// ZoomImage zooms in quarter steps; drags pan.
func ZoomImage(c *ui.Context, label string, image *ui.Bitmap, zoom *float64, o ZoomImageOptions) ui.Element {
	pictures("ZoomImage", label, []*ui.Bitmap{image})
	if o.Max == 0 {
		o.Max = 4
	}
	if !(o.Max > 1) || math.IsInf(o.Max, 0) || zoom == nil || !(*zoom >= 1 && *zoom <= o.Max) {
		panic(fmt.Sprintf("media: ZoomImage %q needs a finite max above 1 and a zoom from 1 to max", label))
	}
	t := c.Theme()
	size("ZoomImage", label, &o.Height, t.Space(64))
	def(&o.In, "Zoom in")
	def(&o.Out, "Zoom out")
	const step = 0.25
	view := ui.Box(c).FillWidth().Height(o.Height).Radius(t.Radius).Clip().Background(t.Surface).Focusable().Label(label)
	b := view.LayoutBounds()
	kit.Settle(c, view, []float32{b.W, b.H}, func() []float32 { r := view.LayoutBounds(); return []float32{r.W, r.H} })
	pan := ui.Local(view, "pan", func() [2]float32 { return [2]float32{} })
	spin := ui.Local(view, "spin", func() float32 { return 0 })
	moved := ui.Local(view, "moved", func() bool { return false })
	view.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll || ev.DY < 0 && *zoom >= o.Max || ev.DY > 0 && *zoom <= 1 {
			return false
		}
		*spin += ev.DY
		return true
	})
	clamp := func() {
		lx, ly := float32(*zoom-1)*b.W/2, float32(*zoom-1)*b.H/2
		pan[0], pan[1] = max(-lx, min(lx, pan[0])), max(-ly, min(ly, pan[1]))
	}
	set := func(z float64) {
		z = max(1, min(o.Max, z))
		if z == *zoom {
			return
		}
		*zoom = z
		clamp()
		view.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("ZoomImage", "zoomed", "zoom", z)
	}
	nudge := func(dx, dy float32) {
		if *zoom > 1 {
			pan[0], pan[1] = pan[0]+dx, pan[1]+dy
			clamp()
			sound.Play(sound.Tick)
		}
	}
	notch := t.Space(10)
	for ; *spin <= -notch; *spin += notch {
		set(*zoom + step)
	}
	for ; *spin >= notch; *spin -= notch {
		set(*zoom - step)
	}
	keys := t.Space(6)
	switch {
	case view.Shortcut(0, ui.KeyEqual), view.Shortcut(ui.Shift, ui.KeyEqual):
		set(*zoom + step)
	case view.Shortcut(0, ui.KeyMinus):
		set(*zoom - step)
	case view.Shortcut(0, ui.Key0):
		set(1)
	case view.Shortcut(0, ui.KeyLeft):
		nudge(keys, 0)
	case view.Shortcut(0, ui.KeyRight):
		nudge(-keys, 0)
	case view.Shortcut(0, ui.KeyUp):
		nudge(0, keys)
	case view.Shortcut(0, ui.KeyDown):
		nudge(0, -keys)
	}
	if dx, dy, held := view.Dragged(); held && *zoom > 1 {
		pan[0], pan[1] = pan[0]+dx, pan[1]+dy
		clamp()
		*moved = *moved || dx != 0 || dy != 0
	} else if !held && *moved {
		*moved = false
		sound.Play(sound.Tick)
		kit.Log("ZoomImage", "panned")
	}
	return view.Children(func() {
		z := view.AnimateTo("z", float32(*zoom), ui.Snappy)
		ui.Image(c, image).Fill().Scale(z, z).
			Translate(view.AnimateTo("px", pan[0], ui.Snappy), view.AnimateTo("py", pan[1], ui.Snappy))
		ui.Row(c).Absolute().End(t.Space(2)).Bottom(t.Space(2)).Gap(t.Space(1)).Children(func() {
			for _, k := range []struct {
				name, glyph string
				by          float64
			}{{o.Out, "zoom-out", -step}, {o.In, "zoom-in", step}} {
				c.Key(k.glyph)
				a := round(c, k.name, k.glyph, t.Space(8))
				a.Background(a.AnimateColor("plate", t.Background.Alpha(0.75+0.25*unit(a.Hovered())), ui.Snappy))
				if a.Clicked() {
					set(*zoom + k.by)
				}
			}
		})
	})
}
