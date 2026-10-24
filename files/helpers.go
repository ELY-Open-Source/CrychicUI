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
	// Path is the folder holding it, when shown.
	Path string
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
	case ".txt", ".md", ".json", ".log", ".pdf", ".doc", ".docx":
		return "file-text"
	case ".csv", ".xlsx", ".numbers":
		return "file-spreadsheet"
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

// kindInk is the theme color of a file's kind.
func kindInk(t *ui.Theme, f FileItem) ui.Color {
	switch itemGlyph(f) {
	case "folder", "file-archive":
		return t.Warning
	case "file-code":
		return t.Accent
	case "image", "file-spreadsheet":
		return t.Success
	case "file-audio", "file-video":
		return t.Danger
	}
	if strings.EqualFold(filepath.Ext(f.Name), ".pdf") {
		return t.Danger
	}
	return t.TextMuted
}

// fileTile is a tinted square holding a file's icon.
func fileTile(c *ui.Context, f FileItem, side float32) ui.Element {
	t := c.Theme()
	ink := kindInk(t, f)
	tile := ui.Box(c).Size(side, side).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.14)).TextColor(ink)
	return tile.Children(func() {
		ui.Icon(c, icon.Must(itemGlyph(f))).Size(side/2, side/2).ContentTransition(ui.ContentTransition{})
	})
}

// shortDate is the default date format of file rows.
func shortDate(t time.Time) string { return t.Format("Jan 2, 2006") }

// rowExit fades a row out; kit.Stagger brings rows in.
var rowExit = ui.ElementTransition{Duration: 160 * time.Millisecond, Exit: &ui.Motion{}}

// fileKey is a file's identity: its folder and name.
func fileKey(f FileItem) string { return fmt.Sprintf("%q %q", f.Path, f.Name) }

// unique panics when two files share a folder and name.
func unique(component, label string, list []FileItem) {
	seen := map[string]bool{}
	for _, f := range list {
		if seen[fileKey(f)] {
			panic(fmt.Sprintf("files: %s %q has two files named %q in %q", component, label, f.Name, f.Path))
		}
		seen[fileKey(f)] = true
	}
}

func sizeText(c *ui.Context, size int64) ui.Element {
	return typography.Bytes(c, size).FontSize(c.Theme().Rem(0.8)).TextColor(c.Theme().TextMuted)
}

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
