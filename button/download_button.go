package button

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// progressButton draws a controlled progress fill and completion.
func progressButton(c *ui.Context, label string, progress float64) ui.Element {
	if math.IsNaN(progress) || math.IsInf(progress, 0) || progress < 0 || progress > 1 {
		panic("button: progress must be between zero and one")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface).Clip()
	done := ui.Local(root, "done", func() bool { return false })
	if progress < 1 {
		*done = false
	} else if !*done {
		*done = true
		sound.Play(sound.Success)
	}
	shown := root.AnimateTo("progress", float32(progress), ui.Snappy)
	root.Draw(func(p *ui.Painter, r ui.Rect) {
		if shown > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * shown, H: r.H}, t.Accent.Alpha(0.2), 0)
		}
	})
	return root.Children(func() {
		if progress >= 1 {
			ui.Icon(c, icon.Must("check")).Size(16, 16).ContentTransition(ui.ContentTransition{})
		} else {
			ui.Text(c, label).SingleLine()
		}
	})
}

// DownloadButton fills with progress and checks when done.
func DownloadButton(c *ui.Context, label string, progress float64) ui.Element {
	return progressButton(c, label, progress)
}
