package media

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Framing is zoom from 1, pan within ±1.
type Framing struct{ Zoom, X, Y float64 }

// AvatarUploadOptions tune an AvatarUpload.
type AvatarUploadOptions struct {
	// Max zoom; zero gives 3.
	Max float64
	// Size of the crop view; zero gives the default.
	Size            float32
	Zoom, Pan, Save string
	// Format writes the zoom; "%.1f×" when empty.
	Format string
}

// AvatarUpload frames an image in a circle; save plays Success.
func AvatarUpload(c *ui.Context, label string, image *ui.Bitmap, f *Framing, o AvatarUploadOptions) (root, save ui.Element) {
	pictures("AvatarUpload", label, []*ui.Bitmap{image})
	if o.Max == 0 {
		o.Max = 3
	}
	if !(o.Max > 1) || math.IsInf(o.Max, 0) || f == nil || !(f.Zoom >= 1 && f.Zoom <= o.Max) || !(f.X >= -1 && f.X <= 1) ||
		!(f.Y >= -1 && f.Y <= 1) {
		panic(fmt.Sprintf("media: AvatarUpload %q needs a max above 1 and a framing within it", label))
	}
	def(&o.Zoom, "Zoom")
	def(&o.Pan, "Position")
	def(&o.Save, "Save")
	def(&o.Format, "%.1f×")
	t := c.Theme()
	size("AvatarUpload", label, &o.Size, t.Space(56))
	sz := o.Size
	root = ui.Column(c).Gap(t.Space(3)).AlignItems(ui.Center).Label(label)
	moved := ui.Local(root, "moved", func() bool { return false })
	root.Children(func() {
		view := kit.Arrive(ui.Box(c).Size(sz, sz).Radius(t.Radius).Clip().Focusable().Label(o.Pan), ui.Pose{Fade: 1, Scale: 0.96}, 0)
		set := func(x, y float64) bool {
			x, y = max(-1, min(1, x)), max(-1, min(1, y))
			if x == f.X && y == f.Y {
				return false
			}
			f.X, f.Y = x, y
			root.SetChanged()
			return true
		}
		w, h := image.Size()
		aw, ah := float64(w)/float64(min(w, h)), float64(h)/float64(min(w, h))
		sx, sy := float64(sz)*(f.Zoom*aw-1)/2, float64(sz)*(f.Zoom*ah-1)/2
		if dx, dy, held := view.Dragged(); held {
			x, y := f.X, f.Y
			if sx > 0 {
				x += float64(dx) / sx
			}
			if sy > 0 {
				y += float64(dy) / sy
			}
			if set(x, y) {
				*moved = true
			}
		} else if *moved {
			*moved = false
			sound.Play(sound.Tick)
			kit.Log("AvatarUpload", "panned", "x", f.X, "y", f.Y)
		}
		for _, k := range []struct {
			key    ui.Key
			dx, dy float64
		}{{ui.KeyLeft, -0.1, 0}, {ui.KeyRight, 0.1, 0}, {ui.KeyUp, 0, -0.1}, {ui.KeyDown, 0, 0.1}} {
			if view.Shortcut(0, k.key) && set(f.X+k.dx, f.Y+k.dy) {
				sound.Play(sound.Tick)
				kit.Log("AvatarUpload", "panned", "x", f.X, "y", f.Y)
			}
		}
		z := view.AnimateTo("zoom", float32(f.Zoom), ui.Snappy)
		px, py := float32(f.X), float32(f.Y)
		view.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			iw, ih := r.W*z*float32(aw), r.H*z*float32(ah)
			pic := ui.Rect{X: cx - iw/2 + px*(iw-r.W)/2, Y: cy - ih/2 + py*(ih-r.H)/2, W: iw, H: ih}
			p.Image(image, pic, ui.FillBox)
			p.Fill(r, kit.Shade(0.5), 0)
			d := r.W - t.Space(4)
			p.Clip(ui.Rect{X: cx - d/2, Y: cy - d/2, W: d, H: d}, d/2, func() { p.Image(image, pic, ui.FillBox) })
			var ring ui.Path
			ring.Circle(cx, cy, d/2)
			p.StrokePath(&ring, t.Space(0.5), t.Background)
		}))
		ui.Box(c).Width(sz).Children(func() {
			if choice.ValueSlider(c, o.Zoom, &f.Zoom, 1, o.Max, choice.ValueSliderOptions{Step: 0.1, Format: o.Format}).FillWidth().Changed() {
				root.SetChanged()
			}
		})
		save = kit.Press(ui.ButtonBase(c).Label(o.Save).Padding(t.Space(1.5), t.Space(3.5)).Radius(t.Radius)).TextColor(t.AccentText)
		bg := t.Accent
		if save.Pressed() {
			bg = t.AccentPressed
		} else if save.Hovered() {
			bg = t.AccentHover
		}
		save.Background(save.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, o.Save).SingleLine() })
		if save.Clicked() {
			sound.Play(sound.Success)
			kit.Log("AvatarUpload", "saved", "zoom", f.Zoom)
		}
	})
	return root, save
}
