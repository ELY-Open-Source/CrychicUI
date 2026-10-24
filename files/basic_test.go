package files

import (
	"bytes"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestFileList(t *testing.T) {
	list := sampleFiles()
	selected, sort, changed := -1, SortName, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if FileList(c, "Files", list, &selected, &sort, FileListOptions{}).Changed() {
			changed++
		}
	}), 600, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Size")
	if sort != SortSize || !strings.Contains(strings.Join(tt.Texts(), " "), "z.txt") {
		t.Fatalf("sort %v texts %q", sort, tt.Texts())
	}
	texts := tt.Texts()
	zi, ai := slices.Index(texts, "z.txt"), slices.Index(texts, "a.go")
	if zi < 0 || ai < 0 || zi > ai {
		t.Errorf("size order %q", texts)
	}
	click(t, tt, "a.go")
	if selected != 1 || changed != 2 {
		t.Errorf("selected %d changed %d", selected, changed)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestFileGrid(t *testing.T) {
	list := sampleFiles()
	selected := 0
	tt := ui.NewTester(padded(func(c *ui.Context) { FileGrid(c, "Grid", list, &selected, FileGridOptions{}) }), 500, 300)
	settles(t, tt)
	beforeImage := bytes.Clone(tt.Image().Pix)
	selected = 1
	tt.Frame()
	settles(t, tt)
	if bytes.Equal(beforeImage, tt.Image().Pix) {
		t.Error("selection has no visible change")
	}
	before := baseline()
	click(t, tt, "a.go")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before)
}

func TestFileIcon(t *testing.T) {
	if glyph("main.go") != "file-code" || glyph("photo.png") != "image" || glyph("unknown") != "file" {
		t.Error("wrong extension icons")
	}
	name := "main.go"
	tt := ui.NewTester(padded(func(c *ui.Context) { FileIcon(c, name, FileIconOptions{}) }), 100, 100)
	settles(t, tt)
	if !tt.HasText("GO") {
		t.Errorf("texts %q", tt.Texts())
	}
	name = "photo.png"
	tt.Frame()
	if !moves(tt, 100*time.Millisecond) || !tt.HasText("PNG") {
		t.Errorf("icon does not swap: %q", tt.Texts())
	}
	panics(t, func() { renderFileIcon("") })
}

func renderFileIcon(name string) {
	ui.Render(func(c *ui.Context) { FileIcon(c, name, FileIconOptions{}) }, 100, 100, 1)
}

func TestDropzone(t *testing.T) {
	paths := []string{}
	tt := ui.NewTester(padded(func(c *ui.Context) { Dropzone(c, "Upload", &paths, DropzoneOptions{}) }), 500, 250)
	settles(t, tt)
	if !tt.HasText("Drop files here") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	called := false
	ui.Render(func(c *ui.Context) {
		root := ui.Box(c)
		if !called {
			appendDrop(root, &paths, []string{"/tmp/a", "/tmp/b"})
			called = true
		}
	}, 100, 100, 1)
	if len(paths) != 2 || paths[1] != "/tmp/b" {
		t.Errorf("paths %q", paths)
	}
	cues(t, before, sound.Drop)
}

func TestUploadList(t *testing.T) {
	uploads := []UploadItem{{Name: "a.txt", Size: 1000, Progress: 0.4}, {Name: "b.go", Size: 2000, Progress: 0.8}, {Name: "c.md", Size: 10, Progress: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) { UploadList(c, "Uploads", &uploads, UploadListOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("40%") || !tt.HasText("400 B") || !tt.HasText("Done c.md") || tt.HasText("Cancel c.md") {
		t.Errorf("texts %q", tt.Texts())
	}
	before := baseline()
	click(t, tt, "Cancel b.go")
	if len(uploads) != 2 || uploads[0].Name != "a.txt" || uploads[1].Name != "c.md" {
		t.Errorf("uploads %v", uploads)
	}
	cues(t, before, sound.Close)
	settles(t, tt)
	before = baseline()
	click(t, tt, "Cancel a.txt")
	uploads = nil
	tt.Frame()
	settles(t, tt)
	if len(uploads) != 0 || !tt.HasText("No uploads") {
		t.Errorf("empty uploads %v texts %q", uploads, tt.Texts())
	}
	cues(t, before, sound.Close)
}
