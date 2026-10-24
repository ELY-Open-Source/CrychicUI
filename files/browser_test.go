package files

import (
	"image"
	"testing"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

func TestStorage(t *testing.T) {
	parts := []StoragePart{{Name: "Images", Bytes: 40}, {Name: "Code", Bytes: 10}}
	tt := ui.NewTester(padded(func(c *ui.Context) { Storage(c, "Storage", parts, 100, StorageOptions{}) }), 500, 300)
	settles(t, tt)
	if !tt.HasText("50 B") || !tt.HasText("of") || !tt.HasText("Images") || !tt.HasText("Free") || !tt.HasText("50%") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestSyncRow(t *testing.T) {
	item := sampleFiles()[0]
	tt := ui.NewTester(padded(func(c *ui.Context) { SyncRow(c, item, Syncing, SyncRowOptions{}) }), 500, 150)
	tt.Frame()
	if !tt.FramePending() || !tt.HasText("Syncing") {
		t.Errorf("pending %v texts %q", tt.FramePending(), tt.Texts())
	}
	done := ui.NewTester(padded(func(c *ui.Context) { SyncRow(c, item, Synced, SyncRowOptions{}) }), 500, 150)
	settles(t, done)
}

func TestDownloadRow(t *testing.T) {
	opened := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		DownloadRow(c, sampleFiles()[0], 1, 20, &opened, DownloadRowOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Open")
	if !opened {
		t.Error("download did not open")
	}
	cues(t, before, sound.Success)
}

func TestZipList(t *testing.T) {
	opened := -1
	entries := []ZipEntry{{Name: "a.txt", Size: 100, Compressed: 40}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		ZipList(c, "Archive", entries, &opened, ZipListOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "a.txt")
	settles(t, tt)
	if opened != 0 || !tt.HasText("Compressed") || !tt.HasText("−60%") {
		t.Errorf("opened %d texts %q", opened, tt.Texts())
	}
	cues(t, before, sound.Toggle)
}

func TestImageFiles(t *testing.T) {
	selected := -1
	bitmap := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	images := []ImageFile{{Item: sampleFiles()[0], Image: bitmap}, {Item: sampleFiles()[1], Image: bitmap}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ImageFiles(c, "Images", images, &selected) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "a.go")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select)
}

func TestCompress(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) { Compress(c, "Compress", 100, 50, CompressOptions{}) }), 500, 180)
	settles(t, tt)
	if !tt.HasText("50 B smaller") {
		t.Errorf("texts %q", tt.Texts())
	}
}

func TestCompressGrowth(t *testing.T) {
	growth := ui.NewTester(padded(func(c *ui.Context) { Compress(c, "Grow", 1, 9, CompressOptions{}) }), 500, 180)
	settles(t, growth)
	empty := ui.NewTester(padded(func(c *ui.Context) { Compress(c, "Empty", 0, 0, CompressOptions{}) }), 500, 180)
	settles(t, empty)
}

func TestCompressCopy(t *testing.T) {
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Compress(c, "Grow", 1, 9, CompressOptions{Larger: "%s more"})
	}), 500, 180)
	if !tt.HasText("8 B more") {
		t.Errorf("custom copy %q", tt.Texts())
	}
}

func TestFileSearch(t *testing.T) {
	query, kind := "", ""
	selected := -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FileSearch(c, "Search", sampleFiles(), &query, &kind, &selected, FileSearchOptions{})
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, ".go")
	if kind != ".go" || tt.HasText("z.txt") {
		t.Errorf("kind %q texts %q", kind, tt.Texts())
	}
	click(t, tt, "a.go")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Select, sound.Select)
}

func TestFileSearchNoopFilter(t *testing.T) {
	query, kind := "", ""
	selected, changed := -1, 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		if FileSearch(c, "Search", sampleFiles(), &query, &kind, &selected, FileSearchOptions{}).Changed() {
			changed++
		}
	}), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "All")
	if changed != 0 {
		t.Errorf("active All reported %d changes", changed)
	}
	cues(t, before)
}

func TestColumnBrowser(t *testing.T) {
	selected := []int{0, 0}
	columns := [][]string{{"home", "docs"}, {"old"}}
	tt := ui.NewTester(padded(func(c *ui.Context) { ColumnBrowser(c, "Browse", columns, &selected) }), 500, 250)
	settles(t, tt)
	before := baseline()
	click(t, tt, "docs")
	if selected[0] != 1 || selected[1] != -1 {
		t.Errorf("selected %v", selected)
	}
	cues(t, before, sound.Select)
}

func TestFileCrumbs(t *testing.T) {
	at := 0
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FileCrumbs(c, "Path", []string{"home", "docs"}, &at, FileCrumbsOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "docs")
	if err := tt.ChooseMenuItem("Open"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if at != 1 {
		t.Errorf("at %d", at)
	}
	cues(t, before, sound.Open)
}

func TestQuotaWarning(t *testing.T) {
	upgraded := false
	tt := ui.NewTester(padded(func(c *ui.Context) {
		QuotaWarning(c, 90, 100, &upgraded, QuotaWarningOptions{})
	}), 500, 200)
	settles(t, tt)
	click(t, tt, "Upgrade")
	if !upgraded {
		t.Error("quota was not upgraded")
	}
	settles(t, tt)
	if !tt.HasText("Plan upgraded") || tt.HasText("Upgrade") {
		t.Errorf("texts %q", tt.Texts())
	}
	if recent := sound.Recent(); recent[len(recent)-1] != sound.Warning {
		t.Errorf("cues %v", recent)
	}
}

func TestStarred(t *testing.T) {
	starred := []bool{false, false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Starred(c, "Starred", sampleFiles(), &starred, StarredOptions{})
	}), 500, 220)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Star a.go")
	if !starred[1] {
		t.Errorf("starred %v", starred)
	}
	cues(t, before, sound.Toggle)
}

func TestConflict(t *testing.T) {
	choice := ConflictNone
	mine, theirs := sampleFiles()[0], sampleFiles()[1]
	tt := ui.NewTester(padded(func(c *ui.Context) { Conflict(c, "Conflict", mine, theirs, &choice, ConflictOptions{}) }), 500, 300)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Keep both")
	if choice != KeepBoth {
		t.Errorf("choice %v", choice)
	}
	cues(t, before, sound.Select)
}

func TestSelectBar(t *testing.T) {
	action := SelectNone
	tt := ui.NewTester(padded(func(c *ui.Context) { SelectBar(c, 3, &action, SelectBarOptions{}) }), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Delete")
	if action != SelectDelete {
		t.Errorf("action %v", action)
	}
	cues(t, before, sound.Tap)
}

func TestFileTags(t *testing.T) {
	selected := []bool{false, false}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FileTags(c, "Tags", []string{"work", "urgent"}, &selected)
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "urgent")
	if !selected[1] {
		t.Errorf("selected %v", selected)
	}
	cues(t, before, sound.Toggle)
}

func TestUploadButtonRow(t *testing.T) {
	action := UploadNone
	tt := ui.NewTester(padded(func(c *ui.Context) {
		UploadButtonRow(c, &action, UploadButtonRowOptions{})
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Upload")
	if action != UploadFiles {
		t.Errorf("action %v", action)
	}
	click(t, tt, "New")
	if err := tt.ChooseMenuItem("New folder"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if action != UploadNewFolder {
		t.Errorf("menu action %v", action)
	}
	cues(t, before, sound.Open, sound.Open)
}

func TestFileTypeFilter(t *testing.T) {
	selected := -1
	types := []TypeCount{{Kind: "Text", Sample: "a.txt", Count: 2}, {Kind: "Go", Sample: "a.go", Count: 1}}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FileTypeFilter(c, "Types", types, &selected)
	}), 500, 180)
	settles(t, tt)
	before := baseline()
	click(t, tt, "Go 1")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Toggle)
}

func TestThumbsCopy(t *testing.T) {
	zoom, selected := 1.0, -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Thumbs(c, "Thumbs", sampleFiles(), &zoom, &selected, ThumbsOptions{Zoom: "Magnification"})
	}), 500, 350)
	if !tt.HasText("Magnification") {
		t.Errorf("custom zoom %q", tt.Texts())
	}
}

func TestThumbs(t *testing.T) {
	zoom, selected := 1.0, -1
	tt := ui.NewTester(padded(func(c *ui.Context) {
		Thumbs(c, "Thumbs", sampleFiles(), &zoom, &selected, ThumbsOptions{})
	}), 500, 350)
	settles(t, tt)
	before := baseline()
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyRight)
	tt.Frame()
	if zoom <= 1 {
		t.Errorf("zoom %v", zoom)
	}
	click(t, tt, "a.go")
	if selected != 1 {
		t.Errorf("selected %d", selected)
	}
	cues(t, before, sound.Tick, sound.Select)
}
