package marketing

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Milestone is one dated step; Icon defaults to a flag.
type Milestone struct{ Year, Title, Text, Icon string }

// HistoryOptions word the timeline.
type HistoryOptions struct{ Title, Text string }

// History is a company timeline.
func History(c *ui.Context, steps []Milestone, o HistoryOptions) ui.Element {
	if len(steps) == 0 {
		panic(fmt.Sprintf("marketing: History %q needs milestones", o.Title))
	}
	for _, s := range steps {
		if strings.TrimSpace(s.Year) == "" || strings.TrimSpace(s.Title) == "" {
			panic(fmt.Sprintf("marketing: History %q milestone %q needs a year and a title", o.Title, s.Title))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		list := ui.Column(c).FillWidth()
		list.Children(func() {
			for i, s := range steps {
				mark, last := s.Icon, i == len(steps)-1
				if mark == "" {
					mark = "flag"
				}
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch)
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Text(c, s.Year).FontSize(t.Rem(0.82)).Bold().TextColor(t.Accent).Width(t.Space(14)).Shrink(0).TextAlign(ui.End).
						Padding(t.Space(2), 0, 0, 0).SingleLine()
					ui.Column(c).Width(t.Space(9)).Shrink(0).AlignItems(ui.Center).Children(func() {
						tile(c, mark, t.Accent, t.Space(9)).Radius(999).Border(1, t.Accent.Alpha(0.3))
						if !last {
							ui.Box(c).Width(1).Grow(1).Background(t.Border)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(1.5), 0, t.Space(5), 0).Children(func() {
						ui.Text(c, s.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().Wrap()
						if s.Text != "" {
							kit.Muted(c, s.Text).FontSize(t.Rem(0.8)).FillWidth().Wrap()
						}
					})
				})
			}
		})
	})
	return root
}
