package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Crop is a box in fractions of the image.
type Crop struct{ X, Y, W, H float64 }

// CropperOptions tune a Cropper.
type CropperOptions struct {
	// Min side as a fraction; zero gives 0.1.
	Min float64
	// Move names the box; Corners its grips, top left first.
	Move    string
	Corners [4]string
}

// Cropper drags and resizes a crop box over an image.
func Cropper(c *ui.Context, label string, image *ui.Bitmap, crop *Crop, o CropperOptions) ui.Element {
	pictures("Cropper", label, []*ui.Bitmap{image})
	if o.Min == 0 {
		o.Min = 0.1
	}
	if !(o.Min > 0 && o.Min <= 1) || crop == nil || !fits(*crop, o.Min) {
		panic(fmt.Sprintf("media: Cropper %q needs a min side in (0, 1] and a crop inside the image", label))
	}
	def(&o.Move, "Crop area")
	for i, d := range [4]string{"Top left corner", "Top right corner", "Bottom left corner", "Bottom right corner"} {
		def(&o.Corners[i], d)
	}
	t := c.Theme()
	w, h := image.Size()
	root := ui.Box(c).FillWidth().AspectRatio(float32(w) / float32(h)).Radius(t.Radius).Clip().Label(label)
	b := root.LayoutBounds()
	kit.Settle(c, root, []float32{b.W, b.H}, func() []float32 { r := root.LayoutBounds(); return []float32{r.W, r.H} })
	moved := ui.Local(root, "moved", func() bool { return false })
	glide := ui.ElementTransition{Animation: &ui.Snappy}
	g := t.Space(3.5)
	var frame ui.Element
	var shades, grips [4]ui.Element
	root.Children(func() {
		ui.Image(c, image).Fill()
		for i := range shades {
			c.Key(i)
			shades[i] = ui.Box(c).Absolute().Background(kit.Shade(0.5)).Transition(glide)
		}
		frame = ui.Box(c).Absolute().Border(t.Space(0.5), t.Background).Focusable().Label(o.Move).Tooltip(o.Move).Transition(glide)
		frame.Children(func() {
			for i := range grips {
				c.Key(i)
				grips[i] = ui.Box(c).Absolute().Size(g, g).Radius(999).Background(t.Background).Border(1, t.Border).
					Label(o.Corners[i])
			}
		})
	})
	set := func(n Crop) {
		if n != *crop {
			*crop = n
			root.SetChanged()
		}
	}
	tick := func(event string) {
		sound.Play(sound.Tick)
		kit.Log("Cropper", event, "x", crop.X, "y", crop.Y, "w", crop.W, "h", crop.H)
	}
	held := false
	for i, grip := range grips {
		dx, dy, on := grip.Dragged()
		if on && b.W > 0 && b.H > 0 {
			held = true
			set(resize(*crop, i, float64(dx/b.W), float64(dy/b.H), o.Min))
			*moved = *moved || dx != 0 || dy != 0
		}
		pose := ui.Pose{}
		if on || grip.Hovered() {
			pose.Scale = 1.25
		}
		grip.AnimatePose(pose, ui.Snappy)
	}
	if dx, dy, on := frame.Dragged(); on && !held && b.W > 0 && b.H > 0 {
		held = true
		set(shift(*crop, float64(dx/b.W), float64(dy/b.H)))
		*moved = *moved || dx != 0 || dy != 0
	}
	if !held && *moved {
		*moved = false
		tick("dragged")
	}
	for _, k := range []struct {
		key    ui.Key
		dx, dy float64
	}{{ui.KeyLeft, -0.01, 0}, {ui.KeyRight, 0.01, 0}, {ui.KeyUp, 0, -0.01}, {ui.KeyDown, 0, 0.01}} {
		if frame.Shortcut(0, k.key) {
			set(shift(*crop, k.dx, k.dy))
			tick("moved")
		}
		if frame.Shortcut(ui.Shift, k.key) {
			set(resize(*crop, 3, k.dx, k.dy, o.Min))
			tick("resized")
		}
	}
	k := *crop
	pct := func(v float64) float32 { return float32(v * 100) }
	frame.LeftPercent(pct(k.X)).TopPercent(pct(k.Y)).WidthPercent(pct(k.W)).HeightPercent(pct(k.H))
	shades[0].Left(0).Right(0).Top(0).HeightPercent(pct(k.Y))
	shades[1].Left(0).Right(0).TopPercent(pct(k.Y + k.H)).Bottom(0)
	shades[2].Left(0).TopPercent(pct(k.Y)).HeightPercent(pct(k.H)).WidthPercent(pct(k.X))
	shades[3].LeftPercent(pct(k.X + k.W)).Right(0).TopPercent(pct(k.Y)).HeightPercent(pct(k.H))
	grips[0].Left(-g / 2).Top(-g / 2)
	grips[1].Right(-g / 2).Top(-g / 2)
	grips[2].Left(-g / 2).Bottom(-g / 2)
	grips[3].Right(-g / 2).Bottom(-g / 2)
	return root
}
