package button

import (
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestRippleButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { RippleButton(c, "Ripple") }, 300, 200)
	before := baseline()
	if err := tt.Click("Ripple"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "RippleButton")
	if !appended(before, sound.Tap) {
		t.Error("missing tap cue")
	}
}

func TestRippleButtonPressOrigin(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(50).Children(func() { RippleButton(c, "Ripple") })
	}, 300, 200)
	r := find(t, tt, "Ripple")
	pressX, pressY := int(r.X+6), int(r.Y+r.H-4)
	releaseX, releaseY := int(r.X+r.W-6), int(r.Y+r.H-4)
	beforePress := tt.Image().RGBAAt(pressX, pressY)
	beforeRelease := tt.Image().RGBAAt(releaseX, releaseY)
	tt.Press(float32(pressX), float32(pressY))
	tt.Move(290, 10)
	tt.Move(float32(releaseX), float32(releaseY))
	tt.Release(float32(releaseX), float32(releaseY))
	tt.Advance(50 * time.Millisecond)
	img := tt.Image()
	if img.RGBAAt(pressX, pressY) == beforePress {
		t.Error("ripple did not start at the press point")
	}
	if img.RGBAAt(releaseX, releaseY) != beforeRelease {
		t.Error("ripple started at the release point")
	}
}

func TestMagneticButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Padding(50).Children(func() { MagneticButton(c, "Magnetic") })
	}, 300, 200)
	before := find(t, tt, "Magnetic")
	r := before
	tt.Move(r.X+2, r.Y+r.H/2)
	settles(t, tt, "MagneticButton toward")
	if moved := find(t, tt, "Magnetic"); moved.X >= before.X {
		t.Errorf("button did not lean: %v to %v", before, moved)
	}
	tt.Move(290, 190)
	settles(t, tt, "MagneticButton back")
	if got := find(t, tt, "Magnetic"); got.X != before.X {
		t.Errorf("button did not return: %v to %v", before, got)
	}
	taps := baseline()
	if err := tt.Click("Magnetic"); err != nil {
		t.Fatal(err)
	}
	if !appended(taps, sound.Tap) {
		t.Error("missing tap cue")
	}
}

func TestShimmerButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { ShimmerButton(c, "Shimmer") }, 300, 200)
	if !tt.FramePending() {
		t.Error("shimmer did not loop")
	}
	tt.Advance(300 * time.Millisecond)
	if !tt.FramePending() {
		t.Error("shimmer stopped")
	}
	taps := baseline()
	if err := tt.Click("Shimmer"); err != nil {
		t.Fatal(err)
	}
	if !appended(taps, sound.Tap) {
		t.Error("missing tap cue")
	}
}

func TestPulseButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { PulseButton(c, "Pulse") }, 300, 200)
	if !tt.FramePending() {
		t.Error("pulse did not loop")
	}
	taps := baseline()
	if err := tt.Click("Pulse"); err != nil {
		t.Fatal(err)
	}
	settles(t, tt, "PulseButton click")
	if !appended(taps, sound.Tap) || tt.FramePending() {
		t.Error("pulse kept running after click")
	}
}

func TestArrowButton(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { ArrowButton(c, "Continue") }, 300, 200)
	r := find(t, tt, "Continue")
	tt.Move(r.X+r.W/2, r.Y+r.H/2)
	settles(t, tt, "ArrowButton hover")
	taps := baseline()
	if err := tt.Click("Continue"); err != nil {
		t.Fatal(err)
	}
	if !appended(taps, sound.Tap) {
		t.Error("missing tap cue")
	}
}
