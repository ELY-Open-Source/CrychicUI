package button

import (
	"math"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestDownloadButton(t *testing.T) {
	progress, clicks := 0.0, 0
	tt := ui.NewTester(func(c *ui.Context) {
		if DownloadButton(c, "Download", progress).Clicked() {
			clicks++
		}
	}, 300, 200)
	before := baseline()
	if err := tt.Click("Download"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 || !appended(before) {
		t.Error("idle click played a completion cue")
	}
	progress = 0.5
	tt.Frame()
	settles(t, tt, "Download progress")
	progress = 1
	tt.Frame()
	settles(t, tt, "Download done")
	if !appended(before, sound.Success) {
		t.Error("completion did not play Success")
	}
	panics(t, "Download progress", func() {
		ui.Render(func(c *ui.Context) { DownloadButton(c, "Download", math.NaN()) }, 10, 10, 1)
	})
}

func TestUploadButton(t *testing.T) {
	path := ""
	progress := 0.0
	tt := ui.NewTester(func(c *ui.Context) { UploadButton(c, "Upload", &path, progress) }, 300, 200)
	before := baseline()
	progress = 1
	tt.Frame()
	settles(t, tt, "Upload done")
	if !appended(before, sound.Success) {
		t.Error("completion did not play Success")
	}
	panics(t, "Upload path", func() { ui.Render(func(c *ui.Context) { UploadButton(c, "Upload", nil, 0) }, 10, 10, 1) })
}

func TestSubmitButton(t *testing.T) {
	state := SubmitIdle
	tt := ui.NewTester(func(c *ui.Context) { SubmitButton(c, "Submit", state) }, 300, 200)
	before := baseline()
	state = SubmitSuccess
	tt.Frame()
	settles(t, tt, "Submit success")
	state = SubmitError
	tt.Frame()
	settles(t, tt, "Submit error")
	if !appended(before, sound.Success, sound.Error) {
		t.Errorf("cues %v", sound.Recent())
	}
	state = SubmitLoading
	tt.Frame()
	if !tt.FramePending() {
		t.Error("loading spinner did not start")
	}
	if err := tt.Click("Submit"); err != nil {
		t.Fatal(err)
	}
	panics(t, "Submit state", func() { ui.Render(func(c *ui.Context) { SubmitButton(c, "Submit", SubmitState(9)) }, 10, 10, 1) })
}
