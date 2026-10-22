package media

import (
	"fmt"
	"time"

	"github.com/egoist/mygo/ui"
)

// Line is a lyric shown from At.
type Line struct {
	At   time.Duration
	Text string
}

// LyricsOptions tune a Lyrics.
type LyricsOptions struct {
	// Height of the view; zero gives the default.
	Height float32
}

// Lyrics lights the line at pos and scrolls to it.
func Lyrics(c *ui.Context, label string, lines []Line, pos time.Duration, o LyricsOptions) ui.Element {
	if label == "" || len(lines) == 0 {
		panic(fmt.Sprintf("media: Lyrics %q needs a label and lines", label))
	}
	for i, l := range lines {
		if l.Text == "" || l.At < 0 || i > 0 && l.At < lines[i-1].At {
			panic(fmt.Sprintf("media: Lyrics %q lines need text and rising times", label))
		}
	}
	t := c.Theme()
	size("Lyrics", label, &o.Height, t.Space(48))
	now := -1
	for i, l := range lines {
		if l.At <= pos {
			now = i
		}
	}
	root := ui.Scroll(c).FillWidth().Height(o.Height).Label(label)
	last := ui.Local(root, "line", func() int { return -2 })
	var lit ui.Element
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4), t.Space(2)).Children(func() {
			for i, l := range lines {
				c.Key(i)
				fg := t.TextMuted
				switch {
				case i == now:
					fg = t.Text
				case i < now:
					fg = t.TextMuted.Alpha(0.55)
				}
				row := ui.Text(c, l.Text).FontSize(t.Rem(1.15)).Bold()
				row.TextColor(row.AnimateColor("fg", fg, ui.Snappy))
				pose := ui.Pose{}
				if i == now {
					pose.Scale = 1.04
					lit = row
				}
				row.TransformOrigin(0, 0.5).AnimatePose(pose, ui.Snappy)
			}
		})
	})
	if now != *last {
		*last = now
		if now >= 0 {
			lit.ScrollIntoView()
		}
	}
	return root
}
