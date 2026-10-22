package media

import (
	"fmt"
	"image"
	"math"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// halves is red on the left, blue on the right.
func halves() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 64, 48))
	for y := range 48 {
		for x := range 64 {
			c := redInk
			if x >= 32 {
				c = blueInk
			}
			img.SetRGBA(x, y, c)
		}
	}
	return ui.NewBitmap(img)
}

func TestImageView(t *testing.T) {
	var img *ui.Bitmap
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() { ImageView(c, "Shot", img, ImageViewOptions{}) })
	}, 400, 400)
	tt.Frame()
	r := find(t, tt, "Shot")
	if at(tt, r.X+r.W/2, r.Y+r.H/2, redInk) {
		t.Error("red shown before the image")
	}
	tt.Advance(400 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("placeholder does not shimmer")
	}
	img = red
	tt.Frame()
	settles(t, tt, "image")
	if !at(tt, r.X+r.W/2, r.Y+r.H/2, redInk) {
		t.Error("image not shown")
	}
	panics(t, "no label", func() { render(func(c *ui.Context) { ImageView(c, "", red, ImageViewOptions{}) }) })
	panics(t, "negative ratio", func() { render(func(c *ui.Context) { ImageView(c, "x", red, ImageViewOptions{Ratio: -1}) }) })
	panics(t, "bad fit", func() { render(func(c *ui.Context) { ImageView(c, "x", red, ImageViewOptions{Fit: 9}) }) })
}

func TestGallery(t *testing.T) {
	open := -1
	imgs := []*ui.Bitmap{red, green, blue}
	tt := ui.NewTester(func(c *ui.Context) { Gallery(c, "Trip", imgs, &open, GalleryOptions{}) }, 600, 500)
	tt.Frame()
	before := baseline()
	click(t, tt, "Photo 2")
	if open != 1 {
		t.Fatalf("open %d, want 1", open)
	}
	settles(t, tt, "lightbox")
	if !at(tt, 300, 260, greenInk) {
		t.Error("lightbox does not show photo 2")
	}
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if open != 2 {
		t.Errorf("open %d after Right twice, want 2", open)
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open != -1 {
		t.Errorf("open %d after Escape, want -1", open)
	}
	cues(t, before, sound.Open, sound.Swipe, sound.Close)
	panics(t, "open out of range", func() {
		n := 3
		render(func(c *ui.Context) { Gallery(c, "x", imgs, &n, GalleryOptions{}) })
	})
	panics(t, "nil image", func() {
		n := -1
		render(func(c *ui.Context) { Gallery(c, "x", []*ui.Bitmap{nil}, &n, GalleryOptions{}) })
	})
}

func TestCarousel(t *testing.T) {
	shown := 0
	tt := ui.NewTester(func(c *ui.Context) {
		Carousel(c, "Tour", 3, &shown, func(i int) { ui.Text(c, fmt.Sprint("Body ", i)) }, CarouselOptions{})
	}, 500, 400)
	settles(t, tt, "carousel")
	view := find(t, tt, "Tour")
	if r, ok := tt.Find("Body 1"); ok && r.W > 1 {
		t.Errorf("slide 1 visible before turning: %v", r)
	}
	before := baseline()
	click(t, tt, "Next")
	if shown != 1 {
		t.Fatalf("at %d after Next, want 1", shown)
	}
	tt.Advance(30 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("slides do not spring")
	}
	settles(t, tt, "carousel")
	if r := find(t, tt, "Body 1"); r.X < view.X-1 || r.X > view.X+view.W/2 {
		t.Errorf("slide 1 at %v, want inside %v", r, view)
	}
	click(t, tt, "Slide 3")
	_ = tt.Click("Next")
	tt.Frame()
	if shown != 2 {
		t.Errorf("at %d, want 2", shown)
	}
	x, y := view.X+view.W/2, view.Y+view.H/2
	tt.Press(x, y)
	for _, dx := range []float32{40, 80, 160} {
		tt.Move(x+dx, y)
		tt.Frame()
	}
	tt.Release(x+160, y)
	tt.Frame()
	if shown != 1 {
		t.Errorf("at %d after a drag right, want 1", shown)
	}
	cues(t, before, sound.Swipe, sound.Swipe, sound.Swipe)
}

func TestCarouselKeys(t *testing.T) {
	shown := 1
	tt := ui.NewTester(func(c *ui.Context) {
		Carousel(c, "Tour", 3, &shown, func(i int) { ui.Text(c, fmt.Sprint("Body ", i)) }, CarouselOptions{})
	}, 500, 400)
	tt.Frame()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if shown != 0 {
		t.Errorf("at %d after Left, want 0", shown)
	}
}

func TestSlideshow(t *testing.T) {
	shown, changed := 0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Slideshow(c, "Show", 3, &shown, func(i int) { ui.Text(c, fmt.Sprint("Frame ", i)) }, SlideshowOptions{Every: time.Second}).Changed() {
			changed++
		}
	}, 500, 400)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("progress bar does not run")
	}
	for i, want := range []int{1, 2, 0} {
		tt.Advance(1050 * time.Millisecond)
		if shown != want {
			t.Fatalf("step %d: at %d, want %d", i, shown, want)
		}
	}
	shown = 2
	tt.Frame()
	tt.Advance(600 * time.Millisecond)
	if shown != 2 {
		t.Errorf("a set slide moved on early: %d", shown)
	}
	if changed != 3 {
		t.Errorf("changed %d times, want 3", changed)
	}
}

func TestThumbStrip(t *testing.T) {
	shown := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ThumbStrip(c, "Strip", []*ui.Bitmap{red, green, blue}, &shown, ThumbStripOptions{})
	}, 500, 400)
	tt.Frame()
	root := find(t, tt, "Strip")
	before := baseline()
	click(t, tt, "Photo 3")
	if shown != 2 {
		t.Fatalf("at %d, want 2", shown)
	}
	settles(t, tt, "strip")
	if !at(tt, root.X+root.W/2, root.Y+60, blueInk) {
		t.Error("main image is not photo 3")
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Frame()
	if shown != 1 {
		t.Errorf("at %d after Left, want 1", shown)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestZoomImage(t *testing.T) {
	z := 1.0
	img := halves()
	tt := ui.NewTester(func(c *ui.Context) { ZoomImage(c, "Map", img, &z, ZoomImageOptions{}) }, 500, 400)
	tt.Frame()
	r := find(t, tt, "Map")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	before := baseline()
	tt.Scroll(cx, cy, 0, -40)
	tt.Frame()
	if z != 1.25 {
		t.Fatalf("zoom %v after a wheel notch, want 1.25", z)
	}
	click(t, tt, "Zoom in")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyEqual)
	tt.Frame()
	if z != 1.75 {
		t.Fatalf("zoom %v, want 1.75", z)
	}
	tt.Key(0, ui.Key0)
	tt.Key(0, ui.Key0)
	tt.Frame()
	if z != 1 {
		t.Errorf("zoom %v after 0, want 1", z)
	}
	tt.Scroll(cx, cy, 0, 40)
	tt.Frame()
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	z = 2
	tt.Frame()
	settles(t, tt, "zoom")
	if !at(tt, cx+12, cy, blueInk) {
		t.Fatal("right of center is not blue at zoom 2")
	}
	before = baseline()
	tt.Press(cx, cy)
	for _, dx := range []float32{20, 40, 60} {
		tt.Move(cx+dx, cy)
		tt.Frame()
	}
	tt.Release(cx+60, cy)
	settles(t, tt, "pan")
	if !at(tt, cx+12, cy, redInk) {
		t.Error("a drag right did not pan the red half under center")
	}
	cues(t, before, sound.Tick)
	panics(t, "max 1", func() { render(func(c *ui.Context) { ZoomImage(c, "x", red, &z, ZoomImageOptions{Max: 1}) }) })
	panics(t, "zoom past max", func() {
		v := 9.0
		render(func(c *ui.Context) { ZoomImage(c, "x", red, &v, ZoomImageOptions{}) })
	})
}

func TestCompare(t *testing.T) {
	split := 0.5
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() { Compare(c, "Edit", red, blue, &split, CompareOptions{}) })
	}, 420, 400)
	settles(t, tt, "compare")
	r := find(t, tt, "Edit")
	y := r.Y + r.H*0.75
	if !at(tt, r.X+r.W*0.2, y, redInk) || !at(tt, r.X+r.W*0.8, y, blueInk) {
		t.Fatal("before should be left, after right")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if math.Abs(split-0.51) > 1e-9 {
		t.Errorf("split %v after Right, want 0.51", split)
	}
	tt.Press(r.X+r.W*0.25, y)
	tt.Frame()
	tt.Release(r.X+r.W*0.25, y)
	settles(t, tt, "compare")
	if math.Abs(split-0.25) > 0.03 {
		t.Errorf("split %v after a press at a quarter", split)
	}
	if !at(tt, r.X+r.W*0.4, y, blueInk) {
		t.Error("after not shown right of a quarter")
	}
	cues(t, before, sound.Tick, sound.Tick)
	panics(t, "split above 1", func() {
		v := 2.0
		render(func(c *ui.Context) { Compare(c, "x", red, blue, &v, CompareOptions{}) })
	})
}

func TestCropper(t *testing.T) {
	crop := Crop{0.2, 0.2, 0.4, 0.4}
	sq := photo(220, 40, 40, 100, 100)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(10).Children(func() { Cropper(c, "Avatar", sq, &crop, CropperOptions{}) })
	}, 320, 340)
	settles(t, tt, "cropper")
	r := find(t, tt, "Avatar")
	if !at(tt, r.X+r.W*0.4, r.Y+r.H*0.4, redInk) || at(tt, r.X+r.W*0.9, r.Y+r.H*0.9, redInk) {
		t.Error("outside the crop should be shaded, inside clear")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Key(ui.Shift, ui.KeyRight)
	tt.Frame()
	if math.Abs(crop.X-0.21) > 1e-9 || math.Abs(crop.W-0.41) > 1e-9 {
		t.Errorf("crop %+v after Right and Shift+Right", crop)
	}
	f := find(t, tt, "Crop area")
	x, y := f.X+f.W/2, f.Y+f.H/2
	tt.Press(x, y)
	for _, dx := range []float32{10, 20, 30} {
		tt.Move(x+dx, y)
		tt.Frame()
	}
	tt.Release(x+30, y)
	tt.Frame()
	if want := 0.21 + 30/float64(r.W); math.Abs(crop.X-want) > 0.01 {
		t.Errorf("x %v after a 30 DIP drag, want %v", crop.X, want)
	}
	g := find(t, tt, "Bottom right corner")
	gx, gy := g.X+g.W/2, g.Y+g.H/2
	tt.Press(gx, gy)
	tt.Move(gx+500, gy+500)
	tt.Frame()
	tt.Release(gx+500, gy+500)
	tt.Frame()
	if crop.X+crop.W > 1+1e-9 || crop.Y+crop.H > 1+1e-9 || math.Abs(crop.X+crop.W-1) > 1e-6 {
		t.Errorf("crop %+v did not stop at the image edge", crop)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Tick)
	panics(t, "crop below min", func() {
		k := Crop{0, 0, 0.05, 0.5}
		render(func(c *ui.Context) { Cropper(c, "x", sq, &k, CropperOptions{}) })
	})
	panics(t, "crop past edge", func() {
		k := Crop{0.8, 0, 0.5, 0.5}
		render(func(c *ui.Context) { Cropper(c, "x", sq, &k, CropperOptions{}) })
	})
}

func TestImageFilters(t *testing.T) {
	f := Original
	tt := ui.NewTester(func(c *ui.Context) { ImageFilters(c, "Looks", red, &f, ImageFiltersOptions{}) }, 520, 420)
	settles(t, tt, "filters")
	r := find(t, tt, "Looks")
	plain := tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+60))
	before := baseline()
	click(t, tt, "Mono")
	if f != Mono {
		t.Fatalf("filter %v, want Mono", f)
	}
	settles(t, tt, "filters")
	p := tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+60))
	if max(p.R, p.G, p.B)-min(p.R, p.G, p.B) > 8 {
		t.Errorf("Mono preview is not gray: %v", p)
	}
	click(t, tt, "Warm")
	settles(t, tt, "filters")
	p = tt.Image().RGBAAt(int(r.X+r.W/2), int(r.Y+60))
	if d := absDiff(p.R, plain.R) + absDiff(p.G, plain.G) + absDiff(p.B, plain.B); d < 15 {
		t.Errorf("Warm preview %v barely differs from the original %v", p, plain)
	}
	cues(t, before, sound.Select, sound.Select)
	panics(t, "two names", func() {
		render(func(c *ui.Context) { ImageFilters(c, "x", red, &f, ImageFiltersOptions{Names: []string{"a", "b"}}) })
	})
	panics(t, "unknown filter", func() {
		g := Filter(9)
		render(func(c *ui.Context) { ImageFilters(c, "x", red, &g, ImageFiltersOptions{}) })
	})
}
