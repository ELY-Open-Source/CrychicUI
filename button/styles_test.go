package button

import (
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTextStyles(t *testing.T) {
	cases := []struct {
		name  string
		cue   sound.Cue
		build func(*ui.Context, string) ui.Element
	}{
		{"Solid", sound.Tap, Solid},
		{"Soft", sound.Tap, Soft},
		{"Outline", sound.Tap, Outline},
		{"Ghost", sound.Tap, Ghost},
		{"Destructive", sound.Delete, Destructive},
		{"Pill", sound.Tap, Pill},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clicks := 0
			tt := ui.NewTester(func(c *ui.Context) {
				if tc.build(c, "Continue").Clicked() {
					clicks++
				}
			}, 300, 200)
			r := find(t, tt, "Continue")
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			settles(t, tt, tc.name+" hover")
			before := sound.Recent()
			if err := tt.Click("Continue"); err != nil {
				t.Fatal(err)
			}
			if clicks != 1 || !appended(before, tc.cue) {
				t.Errorf("clicks %d, cues %v", clicks, sound.Recent())
			}
		})
	}
}

func TestIconButtons(t *testing.T) {
	cases := []struct {
		name  string
		build func(*ui.Context, string) ui.Element
	}{
		{"IconButton", func(c *ui.Context, label string) ui.Element { return IconButton(c, label, "pencil") }},
		{"IconLabel", func(c *ui.Context, label string) ui.Element { return IconLabel(c, label, "pencil") }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) { tc.build(c, "Edit") }, 300, 200)
			r := find(t, tt, "Edit")
			tt.Move(r.X+r.W/2, r.Y+r.H/2)
			settles(t, tt, tc.name)
			taps := sound.Recent()
			if err := tt.Click("Edit"); err != nil {
				t.Fatal(err)
			}
			if !appended(taps, sound.Tap) {
				t.Error("missing tap cue")
			}
		})
	}
	panics(t, "IconButton label", func() { ui.Render(func(c *ui.Context) { IconButton(c, "", "pencil") }, 10, 10, 1) })
}

func TestFab(t *testing.T) {
	clicks := 0
	tt := ui.NewTester(func(c *ui.Context) {
		if Fab(c, "Add", "plus").Clicked() {
			clicks++
		}
	}, 300, 200)
	settles(t, tt, "Fab arrival")
	pops := sound.Recent()
	if err := tt.Click("Add"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 || !appended(pops, sound.Pop) {
		t.Errorf("clicks %d, cues %v", clicks, sound.Recent())
	}
}
