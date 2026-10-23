package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Chapter is a titled start point in a clip.
type Chapter struct {
	Title string
	At    time.Duration
}

// ChaptersOptions tune a Chapters.
type ChaptersOptions struct{}

// Chapters lists chapters; picking one moves pos to its start.
func Chapters(c *ui.Context, label string, chapters []Chapter, pos *time.Duration, o ChaptersOptions) ui.Element {
	if label == "" || len(chapters) == 0 || pos == nil || *pos < 0 {
		panic(fmt.Sprintf("media: Chapters %q needs a label, chapters and a pos not negative", label))
	}
	for i, ch := range chapters {
		if ch.Title == "" || ch.At < 0 || (i > 0 && ch.At <= chapters[i-1].At) {
			panic(fmt.Sprintf("media: Chapters %q chapters need titles and rising starts", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Label(label)
	rows := make([]ui.Element, len(chapters))
	root.Children(func() {
		for i, ch := range chapters {
			c.Key(i)
			rows[i] = kit.Press(ui.ButtonBase(c).Label(ch.Title).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Center).
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
			kit.Arrive(rows[i], ui.Pose{Fade: 1, Y: t.Space(1.5)}, time.Duration(i)*30*time.Millisecond)
		}
	})
	for i, r := range rows {
		if r.Clicked() && *pos != chapters[i].At {
			*pos = chapters[i].At
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("Chapters", "jumped", "index", i)
		}
	}
	now := -1
	for i, ch := range chapters {
		if *pos >= ch.At {
			now = i
		}
	}
	for i, r := range rows {
		on := i == now
		bg, fg := t.SurfaceHover.Alpha(unit(r.Hovered())), t.Text
		if on {
			bg, fg = t.Accent.Alpha(0.1), t.Accent
		}
		r.Checked(on).Background(r.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
			kit.Numeric(c, clock(chapters[i].At)).TextColor(t.TextMuted).Width(t.Space(12)).Shrink(0)
			ui.Text(c, chapters[i].Title).SingleLine().Grow(1).TextColor(r.AnimateColor("fg", fg, ui.Snappy))
		})
	}
	return root
}
