// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"time"
)

type BlameLine struct {
	Author, Text string
	At           time.Time
}

// Blame is lines with an author and age gutter.

// BlameOptions are the copy knobs of Blame.
type BlameOptions struct {
	FormatAge func(age time.Duration) string
}

// Blame is lines with an author and age gutter.
func Blame(c *ui.Context, label string, lines []BlameLine, o BlameOptions) ui.Element {
	if label == "" || len(lines) == 0 {
		panic("code: Blame needs a label and lines")
	}
	if o.FormatAge == nil {
		o.FormatAge = age
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, l := range lines {
			if l.Author == "" || l.At.IsZero() {
				panic("code: Blame needs an author and time in every line")
			}
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(1))
			kit.Arrive(row, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(i*40)*time.Millisecond)
			row.Background(t.Surface.Alpha(unit(float32(i%2) == 0) * 0.5))
			row.Children(func() {
				kit.Muted(c, l.Author).FontSize(t.Rem(0.8)).Width(t.Space(12)).Shrink(0).SingleLine()
				kit.Muted(c, o.FormatAge(c.Now().Sub(l.At))).FontSize(t.Rem(0.8)).Width(t.Space(7)).Shrink(0).SingleLine()
				ui.Text(c, l.Text).Font("monospace").FontSize(t.Rem(0.85)).SingleLine()
			})
		}
	})
	return root
}
