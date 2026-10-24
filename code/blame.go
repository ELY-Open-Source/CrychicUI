package code

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BlameLine is one line of blame.
type BlameLine struct {
	Author, Text string
	At           time.Time
}

// BlameOptions tune a Blame.
type BlameOptions struct {
	// Lang highlights the code; empty leaves it plain.
	Lang      string
	FormatAge func(age time.Duration) string
}

// Blame shows code beside who changed each line and when.
func Blame(c *ui.Context, label string, lines []BlameLine, o BlameOptions) ui.Element {
	if label == "" || len(lines) == 0 {
		panic(fmt.Sprintf("code: Blame %q needs a label and lines", label))
	}
	oldest, newest := lines[0].At, lines[0].At
	for _, l := range lines {
		if l.Author == "" || l.At.IsZero() {
			panic(fmt.Sprintf("code: Blame %q needs an author and time in every line", label))
		}
		if l.At.Before(oldest) {
			oldest = l.At
		}
		if l.At.After(newest) {
			newest = l.At
		}
	}
	if o.FormatAge == nil {
		o.FormatAge = age
	}
	t := c.Theme()
	span := newest.Sub(oldest)
	h := newHighlighter(c, o.Lang)
	state := uint64(0)
	digits := float32(len(strconv.Itoa(len(lines))))
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, l := range lines {
			c.Key(i)
			var spans []ui.Span
			spans, state = h.line(state, l.Text)
			first := i == 0 || lines[i-1].Author != l.Author || !lines[i-1].At.Equal(l.At)
			heat := float32(1)
			if span > 0 {
				heat = float32(l.At.Sub(oldest)) / float32(span)
			}
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(3), 0, 0)
			if first && i > 0 {
				row.BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
			}
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Box(c).Width(3).FillHeight().Shrink(0).Background(t.Accent.Alpha(0.15 + 0.85*heat))
				ui.Row(c).Width(t.Space(36)).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.75), 0).Children(func() {
					if !first {
						return
					}
					display.Avatar(c, l.Author, nil).Size(t.Space(5), t.Space(5)).FontSize(t.Rem(0.55)).Shrink(0)
					ui.Text(c, l.Author).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
					kit.Muted(c, o.FormatAge(c.Now().Sub(l.At))).FontSize(t.Rem(0.7)).Shrink(0)
				})
				kit.Numeric(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted.Alpha(0.7)).
					Width(t.Space(2 * digits)).Shrink(0).TextAlign(ui.End)
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.8)).LineHeight(1.6).Grow(1).MinWidth(0).SingleLine()
			})
		}
	})
	return root
}
