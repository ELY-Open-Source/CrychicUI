package motion

import (
	"fmt"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

// padded builds content inside room to move.
func padded(content func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(40).Children(func() { content(c) }) }
}

// moves checks measure changes within 600 ms, then loops.
func moves(t *testing.T, tt *ui.Tester, what string, measure func() float32) {
	t.Helper()
	a := measure()
	changed := false
	for range 6 {
		tt.Advance(100 * time.Millisecond)
		changed = changed || measure() != a
	}
	if !changed {
		t.Errorf("%s: did not move", what)
	}
	loops(t, tt, what)
}

// stillUnderReduceMotion checks a loop stops when asked for less motion.
func stillUnderReduceMotion(t *testing.T, tt *ui.Tester, what string) {
	t.Helper()
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	tt.Advance(time.Second)
	if tt.FramePending() {
		t.Errorf("%s: loops under reduced motion", what)
	}
}

func TestBob(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Bob(c, func() { dot(c) }) }), 200, 200)
	moves(t, tt, "bob", func() float32 { return find(t, tt, "dot").Y })
	stillUnderReduceMotion(t, tt, "bob")
}

func TestBreathe(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Breathe(c, func() { dot(c) }) }), 200, 200)
	r := find(t, tt, "dot")
	moves(t, tt, "breathe", func() float32 { return float32(tt.Image().RGBAAt(int(r.X+20), int(r.Y+20)).G) })
	stillUnderReduceMotion(t, tt, "breathe")
}

func TestBeat(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Beat(c, func() { dot(c) }) }), 200, 200)
	moves(t, tt, "beat", func() float32 { return find(t, tt, "dot").W })
	stillUnderReduceMotion(t, tt, "beat")
}

func TestSpinning(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Spinning(c, func() { dot(c) }) }), 200, 200)
	moves(t, tt, "spin", func() float32 { return find(t, tt, "dot").W })
}

func TestOrbit(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		Orbit(c, 3, func(i int) { ui.Text(c, fmt.Sprint("p", i)) }, OrbitOptions{})
	}, 200, 200)
	moves(t, tt, "orbit", func() float32 { return find(t, tt, "p0").X })
	panics(t, "no items", func() { render(func(c *ui.Context) { Orbit(c, 0, func(int) {}, OrbitOptions{}) }) })
	panics(t, "NaN radius", func() {
		render(func(c *ui.Context) { Orbit(c, 1, func(int) {}, OrbitOptions{Radius: float32(nan())}) })
	})
}

// kicks checks run changes start motion; first builds rest.
func kicks(t *testing.T, run *uint64, tt *ui.Tester, what string, measure func() float32) {
	t.Helper()
	if tt.FramePending() {
		t.Errorf("%s: moves on first build", what)
	}
	a := measure()
	*run = 1
	tt.Frame()
	changed := false
	for range 8 {
		tt.Advance(40 * time.Millisecond)
		changed = changed || measure() != a
	}
	if !changed {
		t.Errorf("%s: run change did not move it", what)
	}
	tt.Advance(time.Second)
	if tt.FramePending() || measure() != a {
		t.Errorf("%s: did not come to rest where it started", what)
	}
}

func TestShake(t *testing.T) {
	var run uint64
	tt := ui.NewTester(padded(func(c *ui.Context) { Shake(c, run, func() { dot(c) }) }), 200, 200)
	kicks(t, &run, tt, "shake", func() float32 { return find(t, tt, "dot").X })
}

func TestWiggle(t *testing.T) {
	var run uint64
	tt := ui.NewTester(padded(func(c *ui.Context) { Wiggle(c, run, func() { dot(c) }) }), 200, 200)
	kicks(t, &run, tt, "wiggle", func() float32 { return find(t, tt, "dot").W })
}

func TestJump(t *testing.T) {
	var run uint64
	tt := ui.NewTester(padded(func(c *ui.Context) { Jump(c, run, func() { dot(c) }) }), 200, 200)
	kicks(t, &run, tt, "jump", func() float32 { return find(t, tt, "dot").Y })
}
