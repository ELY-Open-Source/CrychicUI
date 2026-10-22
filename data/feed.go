package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Post is a feed card; ID tells cards apart.
type Post struct{ ID, Title, Text string }

// Feed shows cards newest first; new ones slide in.
func Feed(c *ui.Context, label string, posts []Post) ui.Element {
	if label == "" {
		panic("data: Feed needs a label")
	}
	var ids []string
	for _, p := range posts {
		if p.Title == "" {
			panic(fmt.Sprintf("data: Feed %q posts need titles", label))
		}
		ids = append(ids, p.ID)
	}
	if len(ids) > 0 {
		unique("Feed", label, ids)
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList).Label(label)
	seen := ui.Local(root, "seen", func() map[string]bool { return nil })
	fresh := 0
	for _, id := range ids {
		if !(*seen)[id] {
			fresh++
		}
	}
	if *seen != nil && fresh > 0 {
		sound.Play(sound.Notify)
		kit.Log("Feed", "arrived", "count", fresh)
	}
	*seen = map[string]bool{}
	for _, id := range ids {
		(*seen)[id] = true
	}
	slide := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Y: -t.Space(4)}, Exit: &ui.Motion{Collapse: true}, Position: true}
	return root.Children(func() {
		for _, p := range posts {
			c.Key(p.ID)
			kit.Card(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3)).Transition(slide).Role(ui.RoleGroup).Label(p.Title).Children(func() {
				ui.Text(c, p.Title).Bold()
				if p.Text != "" {
					kit.Muted(c, p.Text)
				}
			})
		}
	})
}
