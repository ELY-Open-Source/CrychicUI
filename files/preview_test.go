package files

import (
	"image"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestFilePreview(t *testing.T) {
	item := FileItem{Name: "note.txt", Size: 4, Modified: time.Now()}
	tt := ui.NewTester(padded(func(c *ui.Context) {
		FilePreview(c, "Preview", item, "text", nil, FilePreviewOptions{})
	}), 500, 350)
	settles(t, tt)
	if !tt.HasText("text") {
		t.Errorf("texts %q", tt.Texts())
	}
	img := ui.NewBitmap(image.NewRGBA(image.Rect(0, 0, 2, 2)))
	item.Name = "photo.png"
	imageTest := ui.NewTester(padded(func(c *ui.Context) {
		FilePreview(c, "Image", item, "", img, FilePreviewOptions{})
	}), 500, 350)
	settles(t, imageTest)
	item.Name = "archive.zip"
	info := ui.NewTester(padded(func(c *ui.Context) {
		FilePreview(c, "Info", item, "", nil, FilePreviewOptions{})
	}), 500, 350)
	settles(t, info)
	if !info.HasText("archive.zip") || !info.HasText("ZIP file") {
		t.Errorf("texts %q", info.Texts())
	}
	item.Dir = true
	item.Name = "photos.png"
	folder := ui.NewTester(padded(func(c *ui.Context) {
		FilePreview(c, "Folder", item, "", nil, FilePreviewOptions{})
	}), 500, 350)
	settles(t, folder)
	if !folder.HasText("photos.png") || !folder.HasText("Folder") {
		t.Errorf("folder preview %q", folder.Texts())
	}
	panics(t, func() {
		ui.Render(func(c *ui.Context) {
			FilePreview(c, "Bad", FileItem{Name: "x.png"}, "", nil, FilePreviewOptions{})
		}, 100, 100, 1)
	})
}
