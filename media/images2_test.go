package media

import (
	"bytes"
	"image"
	"image/color"
	"math"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestPhotoGrid(t *testing.T) {
	imgs := []*ui.Bitmap{red, green, green, green, green, green, green}
	tt := ui.NewTester(padded(func(c *ui.Context) { PhotoGrid(c, "Wall", imgs, PhotoGridOptions{}) }), 480, 700)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("tiles do not arrive")
	}
	settles(t, tt, "grid")
	r := find(t, tt, "Wall")
	if !at(tt, r.X+r.W*0.33, r.Y+r.W*0.33, redInk) {
		t.Error("first tile does not span two columns and two rows")
	}
	if !at(tt, r.X+r.W*0.85, r.Y+r.W*0.15, greenInk) {
		t.Error("second tile not beside the first")
	}
	tt.Move(r.X+r.W*0.33, r.Y+r.W*0.33)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("tile does not lift on hover")
	}
	panics(t, "negative columns", func() { render(func(c *ui.Context) { PhotoGrid(c, "x", imgs, PhotoGridOptions{Columns: -1}) }) })
}

func TestAlbumCard(t *testing.T) {
	playing, changed := false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if AlbumCard(c, "Blue", "Joni", red, &playing, AlbumCardOptions{}).Changed() {
			changed++
		}
	}), 400, 400)
	tt.Frame()
	settles(t, tt, "card")
	p := find(t, tt, "Play")
	hidden := tt.Image().RGBAAt(int(p.X+p.W/2), int(p.Y+p.H/2))
	r := find(t, tt, "Blue")
	tt.Move(r.X+r.W/2, r.Y+r.W/2)
	tt.Frame()
	settles(t, tt, "play rising")
	if shown := tt.Image().RGBAAt(int(p.X+p.W/2), int(p.Y+p.H/2)); shown == hidden {
		t.Error("play does not rise on hover")
	}
	before := baseline()
	click(t, tt, "Play")
	if !playing || changed != 1 {
		t.Errorf("playing %v changed %d after Play", playing, changed)
	}
	cues(t, before, sound.Toggle)
	panics(t, "no title", func() {
		on := false
		render(func(c *ui.Context) { AlbumCard(c, "", "", nil, &on, AlbumCardOptions{}) })
	})
}

func TestCoverFlow(t *testing.T) {
	at0, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CoverFlow(c, "Shelf", []*ui.Bitmap{red, green, blue}, &at0, CoverFlowOptions{}).Changed() {
			changed++
		}
	}), 600, 400)
	tt.Frame()
	settles(t, tt, "covers")
	r := find(t, tt, "Shelf")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	if !at(tt, cx, cy, redInk) {
		t.Fatal("first cover not in front")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	settles(t, tt, "turn")
	if at0 != 1 || !at(tt, cx, cy, greenInk) {
		t.Errorf("at %d after Right, want the second cover in front", at0)
	}
	if !at(tt, cx-100, cy, redInk) || !at(tt, cx+100, cy, blueInk) {
		t.Error("side covers do not show beside the front one")
	}
	tt.Key(0, ui.KeyEnd)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	settles(t, tt, "end")
	tt.ClickAt(cx-66, cy)
	tt.Frame()
	if at0 != 2 {
		t.Errorf("at %d after a click on the front cover's edge, want 2", at0)
	}
	tt.ClickAt(cx-90, cy)
	tt.Frame()
	if at0 != 1 || changed != 3 {
		t.Errorf("at %d changed %d after a click on the left cover", at0, changed)
	}
	cues(t, before, sound.Swipe, sound.Swipe, sound.Swipe)
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if at0 != 0 {
		t.Errorf("at %d after Right in right-to-left, want 0", at0)
	}
	panics(t, "at past covers", func() {
		a := 3
		render(func(c *ui.Context) { CoverFlow(c, "x", []*ui.Bitmap{red}, &a, CoverFlowOptions{}) })
	})
}

func TestCollage(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Collage(c, "Mood", []*ui.Bitmap{red, green, blue}, CollageOptions{}) }), 500, 400)
	tt.Frame()
	if !tt.FramePending() {
		t.Error("images do not arrive")
	}
	settles(t, tt, "collage")
	r := find(t, tt, "Mood")
	if !at(tt, r.X+r.W*0.3, r.Y+r.H/2, redInk) || !at(tt, r.X+r.W*0.85, r.Y+r.H/4, greenInk) ||
		!at(tt, r.X+r.W*0.85, r.Y+r.H*3/4, blueInk) {
		t.Error("collage is not one large image beside a column")
	}
	panics(t, "nil image", func() { render(func(c *ui.Context) { Collage(c, "x", []*ui.Bitmap{nil}, CollageOptions{}) }) })
}

func TestPolaroids(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Polaroids(c, "Prints", []*ui.Bitmap{red, green}, []string{"Beach", "Hills"}, PolaroidsOptions{})
	}), 500, 500)
	tt.Frame()
	settles(t, tt, "prints")
	if !slices.Contains(tt.Texts(), "Beach") || !slices.Contains(tt.Texts(), "Hills") {
		t.Errorf("texts %q, want both captions", tt.Texts())
	}
	r := find(t, tt, "Prints")
	tilted := bytes.Clone(tt.Image().Pix)
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "lift")
	if bytes.Equal(tilted, tt.Image().Pix) {
		t.Error("top print does not lift on hover")
	}
	panics(t, "a caption short", func() {
		render(func(c *ui.Context) { Polaroids(c, "x", []*ui.Bitmap{red, green}, []string{"a"}, PolaroidsOptions{}) })
	})
}

func TestMediaGrid(t *testing.T) {
	items := []MediaItem{{"Dunes", red, 0}, {"Trailer", blue, 42 * time.Second}}
	picked, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if MediaGrid(c, "Library", items, &picked, MediaGridOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	tt.Frame()
	if !slices.Contains(tt.Texts(), "0:42") || !slices.Contains(tt.Texts(), "Photo") {
		t.Errorf("texts %q, want the clip length and the photo badge", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Trailer, Video")
	click(t, tt, "Trailer, Video")
	if picked != 1 || changed != 1 {
		t.Errorf("picked %d changed %d", picked, changed)
	}
	cues(t, before, sound.Tap)
	panics(t, "same name twice", func() {
		p := -1
		render(func(c *ui.Context) {
			MediaGrid(c, "x", []MediaItem{{"a", red, 0}, {"a", blue, 0}}, &p, MediaGridOptions{})
		})
	})
}

func TestImageStack(t *testing.T) {
	top, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if ImageStack(c, "Stack", []*ui.Bitmap{red, green, blue}, &top, ImageStackOptions{}).Changed() {
			changed++
		}
	}), 500, 400)
	tt.Frame()
	settles(t, tt, "stack")
	r := find(t, tt, "Stack")
	cx, cy := r.X+r.W/2, r.Y+r.H/2
	if !at(tt, cx, cy, redInk) {
		t.Fatal("first image not in front")
	}
	before := baseline()
	click(t, tt, "Stack")
	if !tt.FramePending() {
		t.Error("cards do not fan to their new places")
	}
	settles(t, tt, "cycle")
	if top != 1 || changed != 1 || !at(tt, cx, cy, greenInk) {
		t.Errorf("top %d changed %d, want the second image in front", top, changed)
	}
	cues(t, before, sound.Swipe)
	panics(t, "top past images", func() {
		k := 3
		render(func(c *ui.Context) { ImageStack(c, "x", []*ui.Bitmap{red}, &k, ImageStackOptions{}) })
	})
}

func TestAvatarUpload(t *testing.T) {
	f, saved, changed := Framing{Zoom: 2, X: 1}, 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		root, save := AvatarUpload(c, "Avatar", halves(), &f, AvatarUploadOptions{})
		if root.Changed() {
			changed++
		}
		if save.Clicked() {
			saved++
		}
	}), 500, 500)
	tt.Frame()
	settles(t, tt, "avatar")
	v := find(t, tt, "Position")
	cx, cy := v.X+v.W/2, v.Y+v.H/2
	if !at(tt, cx, cy, redInk) {
		t.Error("pan right should show the left half in the circle")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyLeft)
	tt.Key(0, ui.KeyRight)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if f.X != 1 || changed != 2 {
		t.Errorf("x %v changed %d; the edge press must change nothing", f.X, changed)
	}
	f.X = -1
	tt.Frame()
	settles(t, tt, "pan")
	if !at(tt, cx, cy, blueInk) {
		t.Error("pan left should show the right half")
	}
	tt.Press(cx, cy)
	for i := range 6 {
		tt.Move(cx+float32(i+1)*5, cy)
		tt.Frame()
	}
	tt.Release(cx+30, cy)
	tt.Frame()
	want := -1 + 30/(float64(v.W)*(2*64.0/48-1)/2)
	if math.Abs(f.X-want) > 0.05 {
		t.Errorf("x %v after a 30 DIP drag, want %v", f.X, want)
	}
	click(t, tt, "Save")
	if saved != 1 {
		t.Errorf("saved %d", saved)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick, sound.Success)
	panics(t, "max at 1", func() {
		g := Framing{Zoom: 1}
		render(func(c *ui.Context) { AvatarUpload(c, "x", red, &g, AvatarUploadOptions{Max: 1}) })
	})
	panics(t, "zoom past max", func() {
		g := Framing{Zoom: 5}
		render(func(c *ui.Context) { AvatarUpload(c, "x", red, &g, AvatarUploadOptions{}) })
	})
}

func thirds() *ui.Bitmap {
	img := image.NewRGBA(image.Rect(0, 0, 300, 100))
	for y := range 100 {
		for x := range 300 {
			img.SetRGBA(x, y, []color.RGBA{redInk, greenInk, blueInk}[x/100])
		}
	}
	return ui.NewBitmap(img)
}

func TestAvatarUploadWide(t *testing.T) {
	f := Framing{Zoom: 2, X: 1}
	pic := thirds()
	tt := ui.NewTester(padded(func(c *ui.Context) { AvatarUpload(c, "Avatar", pic, &f, AvatarUploadOptions{}) }), 500, 500)
	tt.Frame()
	settles(t, tt, "avatar")
	v := find(t, tt, "Position")
	cx, cy := v.X+v.W/2, v.Y+v.H/2
	for _, k := range []struct {
		x    float64
		ink  color.RGBA
		part string
	}{{1, redInk, "left"}, {0, greenInk, "middle"}, {-1, blueInk, "right"}} {
		f.X = k.x
		tt.Frame()
		settles(t, tt, "pan")
		if !at(tt, cx, cy, k.ink) {
			t.Errorf("x %v does not show the %s third", k.x, k.part)
		}
	}
}

func TestCoverFlowSkipsHidden(t *testing.T) {
	at0, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if CoverFlow(c, "Shelf", []*ui.Bitmap{red, green, blue, red, green}, &at0, CoverFlowOptions{}).Changed() {
			changed++
		}
	}), 600, 400)
	tt.Frame()
	settles(t, tt, "covers")
	r := find(t, tt, "Shelf")
	before := baseline()
	tt.ClickAt(r.X+r.W/2+220, r.Y+r.H/2)
	tt.Frame()
	if at0 != 0 || changed != 0 {
		t.Errorf("at %d changed %d after a click on blank space", at0, changed)
	}
	cues(t, before)
}
