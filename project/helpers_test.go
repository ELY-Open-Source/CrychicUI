package project

import (
	"slices"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func padded(view func(c *ui.Context)) func(c *ui.Context) {
	return func(c *ui.Context) { ui.Column(c).Padding(10).Gap(8).Children(func() { view(c) }) }
}

func settles(t *testing.T, tt *ui.Tester) {
	t.Helper()
	tt.Advance(2 * time.Second)
	if tt.FramePending() {
		t.Error("animation did not settle")
	}
}

func click(t *testing.T, tt *ui.Tester, name string) {
	t.Helper()
	if err := tt.Click(name); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
}

func baseline() []sound.Cue {
	for range 32 {
		sound.Play(sound.Warning)
	}
	return sound.Recent()
}

func cues(t *testing.T, before []sound.Cue, want ...sound.Cue) {
	t.Helper()
	after := sound.Recent()
	var got []sound.Cue
	if len(before) == 32 && !slices.Equal(before, after) {
		i := 0
		for i < len(after) && after[i] == sound.Warning {
			i++
		}
		got = after[i:]
	}
	if !slices.Equal(got, want) {
		t.Errorf("cues %v, want %v", got, want)
	}
}

func testTask() Task {
	return Task{ID: "T-1", Title: "Ship UI", Assignee: "Ana", Due: "Friday", Labels: []string{"UI"}}
}
