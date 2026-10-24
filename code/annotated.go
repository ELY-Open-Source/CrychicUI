package code

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Note is a numbered remark on one line.
type Note struct {
	Line int
	Text string
}

// Annotated shows code with numbered notes beside its lines.
func Annotated(c *ui.Context, lang, code string, notes []Note) ui.Element {
	if code == "" || len(notes) == 0 {
		panic("code: Annotated needs code and notes")
	}
	n := len(linesOf(code))
	byLine := map[int]int{}
	last := 0
	for i, nt := range notes {
		if nt.Line < 1 || nt.Line > n || nt.Line <= last || nt.Text == "" {
			panic(fmt.Sprintf("code: Annotated notes must be ascending lines from 1 to %d with text", n))
		}
		byLine[nt.Line] = i + 1
		last = nt.Line
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).Role(ui.RoleGroup).Label(lang)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		codeBody(c, lang, code, true, func(i int, row ui.Element) {
			if byLine[i+1] > 0 {
				row.Background(t.Accent.Alpha(0.07))
			}
		}, func(i int) {
			margin := ui.Row(c).Width(t.Space(46)).Shrink(0).FillHeight().Gap(t.Space(2)).Padding(0, 0, 0, t.Space(3)).
				BorderWidth(0, 0, 0, 1).BorderColor(t.Border)
			m := byLine[i+1]
			if m == 0 {
				return
			}
			kit.Arrive(margin, ui.Pose{X: t.Space(2), Fade: 1}, time.Duration(m)*80*time.Millisecond).Role(ui.RoleGroup).Label(notes[m-1].Text)
			margin.Children(func() {
				ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText).Children(func() {
					ui.Text(c, strconv.Itoa(m)).FontSize(t.Rem(0.7)).Bold()
				})
				ui.Text(c, notes[m-1].Text).FontSize(t.Rem(0.8)).LineHeight(1.4).Grow(1).MinWidth(0).Wrap()
			})
		})
	})
}
