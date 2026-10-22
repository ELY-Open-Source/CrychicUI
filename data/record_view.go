package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Section is a titled set of key and value pairs.
type Section struct {
	Title  string
	Fields []Pair
}

// RecordView shows one record's fields in sections; it fades in.
func RecordView(c *ui.Context, label string, sections []Section) ui.Element {
	if label == "" || len(sections) == 0 {
		panic(fmt.Sprintf("data: RecordView %q needs sections", label))
	}
	var titles []string
	for _, s := range sections {
		titles = append(titles, s.Title)
	}
	unique("RecordView", label, titles)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(5)).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for i, s := range sections {
			c.Key(s.Title)
			sec := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(s.Title)
			kit.Arrive(sec, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*60*time.Millisecond).Children(func() {
				kit.Muted(c, s.Title).Bold().FontSize(t.Rem(0.8))
				ui.Column(c).FillWidth().Dividers(1, t.Border).Children(func() {
					for k, f := range s.Fields {
						c.Key(k)
						ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(1.5), 0).Children(func() {
							kit.Muted(c, f.Key).Width(t.Space(36)).Shrink(0)
							ui.Text(c, f.Value).Grow(1)
						})
					}
				})
			})
		}
	})
}
