package form

import (
	"image"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestTakeKeepsTheFirstDrop(t *testing.T) {
	var path string
	var files []string
	changed := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root := ui.Box(c)
		dropped := files
		files = nil
		take(root, &path, dropped, "Test")
		if root.Changed() {
			changed++
		}
	}, 100, 100)
	files = []string{"/tmp/a.png", "/tmp/b.png"}
	tt.Frame()
	if path != "/tmp/a.png" || changed == 0 {
		t.Errorf("path %q, changed %d", path, changed)
	}
	cues(t, before, sound.Drop)
}

func TestFileField(t *testing.T) {
	path := ""
	browsed, changed := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		root, browse := FileField(c, "Attachment", &path, FileFieldOptions{})
		if root.Changed() {
			changed++
		}
		if browse.Clicked() {
			browsed++
		}
	}, 400, 100)
	find(t, tt, "Drop a file here")
	if tt.HasText("Remove file") {
		t.Error("remove shown with no file")
	}
	click(t, tt, "Browse")
	path = "/tmp/report.pdf"
	tt.Frame()
	find(t, tt, "report.pdf")
	click(t, tt, "Remove file")
	if path != "" || changed != 1 || browsed != 1 {
		t.Errorf("path %q, changed %d, browsed %d", path, changed, browsed)
	}
	cues(t, before, sound.Open, sound.Delete)
	panics(t, "nil path", func() { render(func(c *ui.Context) { FileField(c, "A", nil, FileFieldOptions{}) }) })
}

func TestImageField(t *testing.T) {
	path := "/tmp/cat.png"
	pic := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 4, 4)))
	replaced := 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		if _, replace := ImageField(c, "Photo", &path, pic, ImageFieldOptions{}); replace.Clicked() {
			replaced++
		}
	}, 400, 120)
	click(t, tt, "Replace")
	click(t, tt, "Remove image")
	if path != "" || replaced != 1 {
		t.Errorf("path %q, replaced %d", path, replaced)
	}
	find(t, tt, "Drop an image")
	click(t, tt, "Choose image")
	if replaced != 2 {
		t.Errorf("the empty field offered no picker: %d", replaced)
	}
	cues(t, before, sound.Open, sound.Delete, sound.Open)
}

func TestBugReport(t *testing.T) {
	v := Bug{Severity: "Low"}
	sent, attached := 0, 0
	before := baseline()
	tt := ui.NewTester(func(c *ui.Context) {
		_, send, attach := BugReport(c, &v, []string{"Low", "High"}, BugReportOptions{})
		if send.Clicked() {
			sent++
		}
		if attach.Clicked() {
			attached++
		}
	}, 500, 600)
	click(t, tt, "Send report")
	typeInto(t, tt, "Title", "Crash on save")
	typeInto(t, tt, "Steps to reproduce", "Open, save")
	click(t, tt, "Browse")
	click(t, tt, "Send report")
	if sent != 1 || attached != 1 {
		t.Errorf("sent %d, attached %d, value %+v", sent, attached, v)
	}
	cues(t, before, sound.Open, sound.Send)
	panics(t, "no severities", func() { render(func(c *ui.Context) { BugReport(c, &v, nil, BugReportOptions{}) }) })
}
