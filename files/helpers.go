// Package files holds file management components.
package files

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// FileItem describes one file or folder.
type FileItem struct {
	Name     string
	Size     int64
	Modified time.Time
	Dir      bool
}

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

func dimension(component string, v *float32, d float32) {
	if *v < 0 || math.IsNaN(float64(*v)) || math.IsInf(float64(*v), 0) {
		panic(fmt.Sprintf("files: %s size must be finite and nonnegative", component))
	}
	if *v == 0 {
		*v = d
	}
}

func validFiles(component, label string, list []FileItem) {
	if label == "" || len(list) == 0 {
		panic(fmt.Sprintf("files: %s needs a label and files", component))
	}
	for _, f := range list {
		validFile(component, f)
	}
}

func validFile(component string, f FileItem) {
	if f.Name == "" || f.Size < 0 {
		panic(fmt.Sprintf("files: %s needs a name and nonnegative size", component))
	}
}

func glyph(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".txt", ".md", ".json", ".log":
		return "file-text"
	case ".go", ".js", ".ts", ".css", ".html":
		return "file-code"
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp":
		return "image"
	case ".zip", ".tar", ".gz":
		return "file-archive"
	case ".mp3", ".wav":
		return "file-audio"
	case ".mp4", ".mov":
		return "file-video"
	default:
		return "file"
	}
}

func itemGlyph(f FileItem) string {
	if f.Dir {
		return "folder"
	}
	return glyph(f.Name)
}

func fileIcon(c *ui.Context, name string) ui.Element {
	return ui.Icon(c, icon.Must(glyph(name))).Size(16, 16).ContentTransition(ui.ContentTransition{})
}

func sizeText(c *ui.Context, size int64) ui.Element {
	return typography.Bytes(c, size).FontSize(c.Theme().Rem(0.8)).TextColor(c.Theme().TextMuted)
}

func fileButton(c *ui.Context, label string, chosen bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(chosen))
	bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
	if chosen {
		bg = t.Accent.Alpha(0.12)
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy))
	return b
}

func date(t time.Time) string { return t.Format("2006-01-02") }

func sizeLabel(size int64) string { return fmt.Sprintf("%d B", size) }

func progressBar(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || math.IsNaN(value) || math.IsInf(value, 0) {
		panic("files: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, color, 999)
	}))
	return bar
}
