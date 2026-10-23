package code

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Note is a numbered remark on one line.
type Note struct {
	Line int
	Text string
}

// Annotated shows code with numbered notes beside lines.
func Annotated(c *ui.Context, lang, code string, notes []Note) ui.Element {
	if code == "" {
		panic("code: Annotated needs code")
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
	root := ui.Column(c).FillWidth().Clip().Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Padding(t.Space(2), 0).Role(ui.RoleGroup).Label(lang)
	return root.Children(func() {
		h := newHighlighter(c, lang)
		state := uint64(0)
		for i, line := range linesOf(code) {
			c.Key(i)
			var spans []ui.Span
			spans, state = h.line(state, line)
			ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(0, t.Space(2)).Children(func() {
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.85)).Grow(1)
				m, ok := byLine[i+1]
				if !ok {
					return
				}
				ui.Row(c).Width(t.Space(44)).Shrink(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					chip := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
					kit.Arrive(chip, ui.Pose{Scale: 0.5, Fade: 1}, 0).Children(func() {
						ui.Text(c, strconv.Itoa(m)).FontSize(t.Rem(0.7)).Bold()
					})
					kit.Muted(c, notes[m-1].Text).FontSize(t.Rem(0.8)).Wrap()
				})
			})
		}
	})
}
