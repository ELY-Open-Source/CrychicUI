// Package fx holds backdrops, particles, borders and text fills.
package fx

import (
	"math"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/shader"
)

// scene is a clipped backdrop that content sits on.
func scene(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: a backdrop needs content")
	}
	t := c.Theme()
	return ui.Box(c).FillWidth().MinHeight(t.Space(40)).Radius(t.Radius).Clip()
}

// phase loops 0 to 1; reduced motion holds 0.
func phase(c *ui.Context, e ui.Element, period time.Duration) float32 {
	if c.Preferences().ReduceMotion {
		return 0
	}
	return e.Loop("crychic.phase", period, ui.Linear)
}

// rnd is a fixed pseudo-random fraction for i and salt.
func rnd(i int, salt uint32) float32 {
	s := uint32(i)*747796405 + 2891336453 + salt*2654435761
	w := ((s >> ((s >> 28) + 4)) ^ s) * 277803737
	return float32((w>>22)^w) / float32(math.MaxUint32+1.0)
}

// fract is the part of v after its floor.
func fract(v float32) float32 { return v - float32(math.Floor(float64(v))) }

func unit(on bool) float32 {
	if on {
		return 1
	}
	return 0
}

// palette is the theme's four signal colors.
func palette(t *ui.Theme) []ui.Color { return []ui.Color{t.Accent, t.Success, t.Warning, t.Danger} }

// at is a gradient stop of col at place.
func at(col ui.Color, place float32) ui.GradientStop { return ui.GradientStop{Color: col, At: place} }

// segment strokes a straight line at any angle.
func segment(p *ui.Painter, x0, y0, x1, y1, width float32, col ui.Color) {
	var l ui.Path
	p.StrokePath(l.MoveTo(x0, y0).LineTo(x1, y1), width, col)
}

// glow paints a soft disc fading out from its center.
func glow(p *ui.Painter, x, y, radius float32, col ui.Color, alpha float32) {
	const rings = 12
	for i := range rings {
		var disc ui.Path
		p.FillPath(disc.Circle(x, y, radius*(1-float32(i)/rings)), col.Alpha(alpha/rings))
	}
}

// follow eases toward the pointer; a fades with hover.
func follow(root ui.Element) (x, y, a float32) {
	px, py, over := root.PointerPosition()
	was := ui.Local(root, "crychic.over", func() bool { return false })
	if over && !*was {
		root.Jump("crychic.x", px)
		root.Jump("crychic.y", py)
	}
	at := ui.Local(root, "crychic.at", func() ui.Point { return ui.Point{} })
	if *was = over; over {
		*at = ui.Point{X: px, Y: py}
	}
	return root.AnimateTo("crychic.x", at.X, ui.Smooth), root.AnimateTo("crychic.y", at.Y, ui.Smooth),
		root.AnimateTo("crychic.a", unit(over), ui.Smooth)
}

// burst plays 0 to 1 once run leaves its start.
func burst(e ui.Element, run uint64, d time.Duration) float32 {
	first := ui.Local(e, "crychic.first", func() uint64 { return run })
	moved := ui.Local(e, "crychic.moved", func() bool { return false })
	if run != *first {
		*moved = true
	}
	if !*moved {
		return 1
	}
	return e.Keyframes("crychic.burst", []ui.Keyframe{{Value: 0}, {At: time.Millisecond, Value: 0},
		{At: d, Value: 1, Ease: ui.Linear}}, ui.KeyframeOptions{Run: run})
}

// tap is a click's place in its box, and when.
type tap struct {
	x, y float32
	at   time.Time
}

// taps cues each click on root, kept for life.
func taps(c *ui.Context, root ui.Element, name string, cue sound.Cue, life time.Duration) []tap {
	list := ui.Local(root, "crychic.taps", func() []tap { return nil })
	now := c.Now()
	if root.Clicked() {
		x, y, _ := root.PointerPosition()
		*list = append(*list, tap{x: x, y: y, at: now})
		sound.Play(cue)
		kit.Log(name, "tap", "count", len(*list))
	}
	*list = slices.DeleteFunc(*list, func(t tap) bool { return now.Sub(t.at) >= life })
	if c.Preferences().ReduceMotion {
		*list = nil
	}
	if len(*list) > 0 {
		c.AnimationFrame()
	}
	return slices.Clone(*list)
}

// paint fills A to B; Flow waves stops by Shift.
var paint = shader.MustCompile(`package main

var A vec4
var B vec4
var Shift float
var Flow float

func Fragment(dst vec4, pos vec2) vec4 {
	x := pos.x / rectSize().x
	f := mix(x, 0.5+0.5*sin(6.2832*(x-Shift)), Flow)
	c := mix(A, B, f)
	return vec4(c.rgb*c.a, c.a)
}
`)

// headline is the bold span the text fills draw.
func headline(c *ui.Context, text string) ui.Span {
	if text == "" {
		panic("fx: text effects need text")
	}
	return ui.Span{Text: text, Size: c.Theme().Rem(1.5), Weight: 700}
}

// filled is a labeled box sized to s.
func filled(c *ui.Context, s ui.Span) ui.Element {
	w, h := c.MeasureText(0, s)
	return ui.Box(c).Size(w, h).Shrink(0).Role(ui.RoleText).Label(s.Text)
}
