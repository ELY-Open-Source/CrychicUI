package media

import (
	"fmt"
	"image/color"
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// inkNear finds c within 3 DIP of x.
func inkNear(tt *ui.Tester, x, y float32, c color.RGBA) bool {
	for dx := float32(-3); dx <= 3; dx++ {
		if at(tt, x+dx, y, c) {
			return true
		}
	}
	return false
}

func TestCaptions(t *testing.T) {
	lines := []Caption{{0, time.Second, "Hello"}, {2 * time.Second, 3 * time.Second, "World"}}
	pos := 500 * time.Millisecond
	tt := ui.NewTester(padded(func(c *ui.Context) { Captions(c, "Subs", lines, pos, CaptionsOptions{}) }), 400, 200)
	tt.Frame()
	if !slices.Contains(tt.Texts(), "Hello") {
		t.Fatalf("texts %q, want Hello", tt.Texts())
	}
	pos = 1500 * time.Millisecond
	tt.Frame()
	settles(t, tt, "gap")
	if slices.Contains(tt.Texts(), "Hello") || slices.Contains(tt.Texts(), "World") {
		t.Errorf("texts %q in the gap", tt.Texts())
	}
	pos = 2500 * time.Millisecond
	tt.Frame()
	if !tt.FramePending() {
		t.Error("caption does not fade in")
	}
	settles(t, tt, "caption")
	if !slices.Contains(tt.Texts(), "World") {
		t.Errorf("texts %q, want World", tt.Texts())
	}
	for name, bad := range map[string][]Caption{
		"overlap": {{0, 2 * time.Second, "a"}, {time.Second, 3 * time.Second, "b"}},
		"empty":   {{0, time.Second, ""}},
		"reverse": {{time.Second, time.Second, "a"}},
	} {
		panics(t, name, func() { render(func(c *ui.Context) { Captions(c, "x", bad, 0, CaptionsOptions{}) }) })
	}
}

func TestChapters(t *testing.T) {
	chs := []Chapter{{"Intro", 0}, {"Verse", 30 * time.Second}, {"Chorus", time.Minute}}
	pos, changed := 45*time.Second, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Chapters(c, "Parts", chs, &pos, ChaptersOptions{}).Changed() {
			changed++
		}
	}), 400, 300)
	tt.Frame()
	before := baseline()
	click(t, tt, "Chorus")
	if pos != time.Minute || changed != 1 {
		t.Errorf("pos %v changed %d after Chorus", pos, changed)
	}
	click(t, tt, "Chorus")
	if changed != 1 {
		t.Error("picking the current start reported a change")
	}
	cues(t, before, sound.Select)
	panics(t, "falling starts", func() {
		p := time.Duration(0)
		render(func(c *ui.Context) { Chapters(c, "x", []Chapter{{"a", time.Second}, {"b", 0}}, &p, ChaptersOptions{}) })
	})
}

func TestStories(t *testing.T) {
	idx, changed := 0, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Stories(c, "Story", 3, &idx, func(i int) { ui.Text(c, fmt.Sprintf("Slide %d", i+1)) }, StoriesOptions{Every: time.Second}).Changed() {
			changed++
		}
	}), 400, 420)
	tt.Frame()
	before := baseline()
	tt.Advance(1100 * time.Millisecond)
	if idx != 1 {
		t.Fatalf("at %d after one slide's time, want 1", idx)
	}
	r := find(t, tt, "Story")
	tt.ClickAt(r.X+r.W*0.8, r.Y+r.H/2)
	tt.Frame()
	tt.ClickAt(r.X+r.W*0.8, r.Y+r.H/2)
	tt.Frame()
	if idx != 2 {
		t.Errorf("at %d after two taps on the right, want 2", idx)
	}
	tt.ClickAt(r.X+r.W*0.1, r.Y+r.H/2)
	tt.Frame()
	if idx != 1 {
		t.Errorf("at %d after a tap on the left, want 1", idx)
	}
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	tt.Advance(3 * time.Second)
	if idx != 2 {
		t.Errorf("at %d, want it to stop on the last slide", idx)
	}
	if changed != 4 {
		t.Errorf("changed %d times, want timer, tap, tap and key", changed)
	}
	cues(t, before, sound.Tick, sound.Tick, sound.Tick)
	full := 0
	tf := ui.NewTester(padded(func(c *ui.Context) {
		Stories(c, "Photos", 2, &full, func(int) { ui.Image(c, red).Fill().Fit(ui.Cover) }, StoriesOptions{})
	}), 400, 420)
	tf.Advance(500 * time.Millisecond)
	if f := find(t, tf, "Photos"); !at(tf, f.X+f.W/2, f.Y+f.H*0.9, redInk) {
		t.Error("slide does not fill its frame")
	}
	panics(t, "no slides", func() {
		a := 0
		render(func(c *ui.Context) { Stories(c, "x", 0, &a, func(int) {}, StoriesOptions{}) })
	})
}

func TestFilmstrip(t *testing.T) {
	var acc color.RGBA
	pos := 5 * time.Second
	frames := []*ui.Bitmap{red, red, red, red}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		acc = rgba(c.Theme().Accent)
		Filmstrip(c, "Reel", frames, pos, 10*time.Second, FilmstripOptions{})
	}), 500, 200)
	tt.Frame()
	settles(t, tt, "film")
	r := find(t, tt, "Reel")
	y := r.Y + r.H/2
	if !inkNear(tt, r.X+r.W/2, y, acc) {
		t.Error("no playhead at the middle")
	}
	pos = 8 * time.Second
	tt.Frame()
	settles(t, tt, "playhead")
	if !inkNear(tt, r.X+r.W*0.8, y, acc) || inkNear(tt, r.X+r.W/2, y, acc) {
		t.Error("playhead did not move to 80%")
	}
	panics(t, "pos past total", func() {
		render(func(c *ui.Context) { Filmstrip(c, "x", frames, 11*time.Second, 10*time.Second, FilmstripOptions{}) })
	})
}

func TestScrubber(t *testing.T) {
	pos, changed := time.Duration(0), 0
	frames := []*ui.Bitmap{red, green, blue}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if Scrubber(c, "Scrub", frames, &pos, 30*time.Second, ScrubberOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	tt.Frame()
	r := find(t, tt, "Scrub")
	shows := func(c color.RGBA) bool {
		for y := r.Y - 100; y < r.Y; y += 2 {
			for x := r.X; x < r.X+r.W; x += 2 {
				if at(tt, x, y, c) {
					return true
				}
			}
		}
		return false
	}
	tt.Move(r.X+r.W*0.85, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "preview")
	if !shows(blueInk) || shows(greenInk) {
		t.Error("hover at 85% should preview the last frame")
	}
	tt.Move(r.X+r.W*0.45, r.Y+r.H/2)
	tt.Frame()
	settles(t, tt, "preview")
	if !shows(greenInk) {
		t.Error("hover at 45% should preview the middle frame")
	}
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if pos != time.Second || changed != 1 {
		t.Errorf("pos %v changed %d after Right", pos, changed)
	}
	cues(t, before, sound.Tick)
	panics(t, "pos past total", func() {
		p := 31 * time.Second
		render(func(c *ui.Context) { Scrubber(c, "x", frames, &p, 30*time.Second, ScrubberOptions{}) })
	})
}

func TestSpeedMenu(t *testing.T) {
	speed, open, changed := 1.0, false, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if SpeedMenu(c, "Speed", []float64{0.5, 1, 1.5, 2}, &speed, &open, SpeedMenuOptions{}).Changed() {
			changed++
		}
	}), 400, 400)
	tt.Frame()
	before := baseline()
	click(t, tt, "Speed: 1×")
	if !open {
		t.Fatal("menu did not open")
	}
	tt.Advance(300 * time.Millisecond)
	click(t, tt, "1.5×")
	if speed != 1.5 || open {
		t.Errorf("speed %v open %v after 1.5×", speed, open)
	}
	click(t, tt, "Speed: 1.5×")
	tt.Advance(300 * time.Millisecond)
	tt.Key(0, ui.KeyDown)
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if speed != 2 {
		t.Errorf("speed %v after Down, Enter from 1.5×; the checked row should take focus", speed)
	}
	if changed != 4 {
		t.Errorf("changed %d, want 4", changed)
	}
	cues(t, before, sound.Open, sound.Select, sound.Open, sound.Select)
	for name, s := range map[string][]float64{"zero": {0, 1}, "twice": {1, 1}} {
		panics(t, name, func() {
			v, o := 1.0, false
			render(func(c *ui.Context) { SpeedMenu(c, "x", s, &v, &o, SpeedMenuOptions{}) })
		})
	}
	panics(t, "speed not listed", func() {
		v, o := 3.0, false
		render(func(c *ui.Context) { SpeedMenu(c, "x", []float64{1, 2}, &v, &o, SpeedMenuOptions{}) })
	})
}

func TestQualityMenu(t *testing.T) {
	at, open := 0, false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QualityMenu(c, "Quality", []string{"Auto", "1080p", "720p"}, &at, &open, QualityMenuOptions{})
	}), 400, 400)
	tt.Frame()
	before := baseline()
	click(t, tt, "Quality: Auto")
	tt.Advance(300 * time.Millisecond)
	click(t, tt, "720p")
	if at != 2 || open {
		t.Errorf("at %d open %v after 720p", at, open)
	}
	click(t, tt, "Quality: 720p")
	click(t, tt, "Quality: 720p")
	if open {
		t.Error("second trigger click did not close")
	}
	cues(t, before, sound.Open, sound.Select, sound.Open, sound.Close)
	panics(t, "twice", func() {
		a, o := 0, false
		render(func(c *ui.Context) { QualityMenu(c, "x", []string{"a", "a"}, &a, &o, QualityMenuOptions{}) })
	})
}

func TestVideoFrame(t *testing.T) {
	pos, playing := time.Duration(0), false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		VideoFrame(c, "Clip", red, &pos, 10*time.Second, &playing, VideoFrameOptions{})
	}), 500, 400)
	tt.Frame()
	settles(t, tt, "frame")
	r := find(t, tt, "Clip")
	cx, cy := r.X+r.W/2, r.Y+r.W*9/32
	if at(tt, cx, cy, redInk) {
		t.Error("no play overlay over the poster")
	}
	before := baseline()
	click(t, tt, "Play")
	if !playing {
		t.Fatal("overlay did not start play")
	}
	tt.Advance(time.Second)
	if d := pos - time.Second; d < -40*time.Millisecond || d > 40*time.Millisecond {
		t.Errorf("pos %v after a second", pos)
	}
	if !at(tt, cx, cy, redInk) {
		t.Error("overlay still shown while playing")
	}
	playing = false
	tt.Frame()
	held := pos
	tt.Advance(time.Second)
	if pos != held {
		t.Errorf("pos moved from %v to %v while paused", held, pos)
	}
	cues(t, before, sound.Toggle)
	panics(t, "no poster", func() {
		p, on := time.Duration(0), false
		render(func(c *ui.Context) { VideoFrame(c, "x", nil, &p, time.Second, &on, VideoFrameOptions{}) })
	})
}
