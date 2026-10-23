// Package media holds image, audio and player components.
package media

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// def fills an empty string with its default.
func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

func abs(v float32) float32 { return max(v, -v) }

// clock writes d as m:ss, or h:mm:ss past an hour.
func clock(d time.Duration) string {
	s := int64(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// pictures panics on a missing label or image.
func pictures(component, label string, images []*ui.Bitmap) {
	if label == "" || len(images) == 0 {
		panic(fmt.Sprintf("media: %s %q needs a label and images", component, label))
	}
	for _, img := range images {
		if img == nil {
			panic(fmt.Sprintf("media: %s %q images must not be nil", component, label))
		}
	}
}

// chosen panics unless at runs from lo to n-1.
func chosen(component, label string, at *int, lo, n int) {
	if at == nil || *at < lo || *at >= n {
		panic(fmt.Sprintf("media: %s %q needs a choice from %d to %d", component, label, lo, n-1))
	}
}

// size panics on a bad size; zero becomes d.
func size(component, label string, v *float32, d float32) {
	if *v < 0 || math.IsNaN(float64(*v)) || math.IsInf(float64(*v), 0) {
		panic(fmt.Sprintf("media: %s %q sizes must be finite and not negative", component, label))
	}
	if *v == 0 {
		*v = d
	}
}

// tracked panics on a track without a title or length.
func tracked(component string, k Track) {
	if k.Title == "" || k.Length <= 0 {
		panic(fmt.Sprintf("media: %s tracks need a title and a length", component))
	}
}

// cover is a track's art, or a music glyph.
func cover(c *ui.Context, art *ui.Bitmap, side float32) ui.Element {
	t := c.Theme()
	b := ui.Box(c).Size(side, side).Shrink(0).Radius(t.Radius).Clip().Background(t.Surface).Center()
	return b.Children(func() {
		if art != nil {
			ui.Image(c, art).Fill().Fit(ui.Cover)
			return
		}
		ui.Icon(c, icon.Must("music")).Size(16, 16).TextColor(t.TextMuted)
	})
}

// round is an icon button; the caller plays its cue.
func round(c *ui.Context, name, glyph string, side float32) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(side, side).Shrink(0).Center().Radius(999))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	return b.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
}

// ends gives back and forth glyphs and keys.
func ends(c *ui.Context) (back, forth string, prev, next ui.Key) {
	back, forth, prev, next = "chevron-left", "chevron-right", ui.KeyLeft, ui.KeyRight
	if c.Direction() == ui.RightToLeft {
		back, forth, prev, next = forth, back, next, prev
	}
	return back, forth, prev, next
}

// released plays Tick after a drag or key move.
func released(s ui.Element, component string, value float64) {
	moved := ui.Local(s, "moved", func() bool { return false })
	*moved = *moved || s.Changed()
	if _, _, held := s.Dragged(); *moved && !held {
		*moved = false
		sound.Play(sound.Tick)
		kit.Log(component, "changed", "value", value)
	}
}

// fits reports whether a crop lies inside the image.
func fits(k Crop, least float64) bool {
	const eps = 1e-9
	return k.X >= -eps && k.Y >= -eps && k.W >= least-eps && k.H >= least-eps && k.X+k.W <= 1+eps && k.Y+k.H <= 1+eps
}

// shift moves a crop, kept inside the image.
func shift(k Crop, dx, dy float64) Crop {
	k.X = max(0, min(1-k.W, k.X+dx))
	k.Y = max(0, min(1-k.H, k.Y+dy))
	return k
}

// resize drags crop corner i, top left first.
func resize(k Crop, i int, dx, dy, least float64) Crop {
	if i%2 == 0 {
		x := max(0, min(k.X+k.W-least, k.X+dx))
		k.W, k.X = k.W+k.X-x, x
	} else {
		k.W = max(least, min(1-k.X, k.W+dx))
	}
	if i < 2 {
		y := max(0, min(k.Y+k.H-least, k.Y+dy))
		k.H, k.Y = k.H+k.Y-y, y
	} else {
		k.H = max(least, min(1-k.Y, k.H+dy))
	}
	return k
}

// look draws an image in filter f.
func look(c *ui.Context, image *ui.Bitmap, f Filter) ui.Element {
	t := c.Theme()
	return ui.Grid(c).Columns(1).Fill().Children(func() {
		img := ui.Image(c, image).ColumnStart(1).RowStart(1).Fill().Fit(ui.Cover)
		var tint ui.Color
		switch f {
		case Mono:
			img.Grayscale()
		case Soft:
			img.Blur(t.Space(0.75))
		case Warm:
			tint = t.Warning.Alpha(0.25)
		case Cool:
			tint = t.Accent.Alpha(0.22)
		case Faded:
			tint = t.Background.Alpha(0.4)
		}
		if f >= Warm {
			ui.Box(c).ColumnStart(1).RowStart(1).Fill().Background(tint)
		}
	})
}

// advance runs pos while *playing; true when play restarts.
func advance(root ui.Element, now time.Time, pos *time.Duration, total time.Duration, playing *bool, component string) bool {
	last := ui.Local(root, "last", func() time.Time { return time.Time{} })
	switch {
	case !*playing:
		*last = time.Time{}
	case last.IsZero():
		if *pos >= total {
			*pos = 0
		}
		*last = now
		return true
	default:
		*pos = min(total, *pos+now.Sub(*last))
		*last = now
		if *pos >= total {
			*playing, *last = false, time.Time{}
			root.SetChanged()
			kit.Log(component, "ended")
		}
	}
	return false
}

// spin gives e's turn, coasting to rest when stopped.
func spin(c *ui.Context, e ui.Element, playing bool, rps float32) float32 {
	now := c.Now()
	turn := ui.Local(e, "turn", func() float64 { return 0 })
	last := ui.Local(e, "spun", func() time.Time { return now })
	speed := e.AnimateTo("speed", rps*unit(playing && !c.Preferences().ReduceMotion), ui.Smooth)
	*turn = math.Mod(*turn+float64(speed)*now.Sub(*last).Seconds(), 1)
	*last = now
	if speed != 0 {
		c.AnimationFrame()
	}
	return float32(*turn)
}

// bands gives n sprung bar heights from 0 to 1.
func bands(c *ui.Context, e ui.Element, n int, playing bool) []float32 {
	live := playing && !c.Preferences().ReduceMotion
	sec := float64(c.Now().UnixMilli()%600000) / 1000
	hs := make([]float32, n)
	for i := range hs {
		f := float64(i) / float64(n)
		aim := 0.06
		switch {
		case live:
			aim = (0.15 + 0.85*math.Abs(math.Sin(sec*(2.1+1.7*f)+float64(i)*0.9))) * (1 - 0.55*f)
		case playing:
			aim = 0.6 * (1 - 0.55*f)
		}
		hs[i] = e.AnimateTo(fmt.Sprint("bar.", i), float32(aim), ui.Snappy)
	}
	if live {
		c.AnimationFrame()
	}
	return hs
}

// white reports whether MIDI note m is a white key.
func white(m int) bool {
	switch m % 12 {
	case 1, 3, 6, 8, 10:
		return false
	}
	return true
}

// pitch is MIDI note m in Hz.
func pitch(m int) float64 { return 440 * math.Pow(2, float64(m-69)/12) }

// named writes MIDI note m as C4 or F#5.
func named(m int) string {
	names := [12]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	return fmt.Sprintf("%s%d", names[m%12], m/12-1)
}

// strike reports a pointer, Enter or Space press.
func strike(e ui.Element) bool {
	return struck(e) || e.Shortcut(0, ui.KeyEnter) || e.Shortcut(0, ui.KeySpace)
}

// flash dips to 1 and back per strike.
func flash(e ui.Element, run int) float32 {
	if run == 0 {
		return 0
	}
	return e.Keyframes(run, []ui.Keyframe{{Value: 1}, {At: 180 * time.Millisecond, Value: 0, Ease: ui.EaseOut}})
}

// struck reports a pointer press on e, once per press.
func struck(e ui.Element) bool {
	hit := ui.Local(e, "struck", func() bool { return false })
	e.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown && ev.Button == 0 {
			*hit = true
		}
		return false
	})
	was := *hit
	*hit = false
	return was
}
