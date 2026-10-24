package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Change is one task changelog entry.
type Change struct {
	Text, Time string
}

// Changes shows a task changelog.
func Changes(c *ui.Context, entries []Change) ui.Element {
	if len(entries) == 0 {
		panic("project: Changes needs entries")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, entry := range entries {
			if entry.Text == "" || entry.Time == "" {
				panic("project: Changes has an invalid entry")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, entry.Text).Grow(1).Wrap()
				kit.Muted(c, entry.Time).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
