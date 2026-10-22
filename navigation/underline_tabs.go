package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnderlineTabs are tabs whose underline slides to the chosen one.
func UnderlineTabs(c *ui.Context, label string, names []string, active *int) ui.Element {
	tabs("UnderlineTabs", label, names, active)
	t := c.Theme()
	shown := *active
	parts := ui.TabsBase(c, active, len(names))
	root := parts.List.Label(label).Gap(t.Space(1))
	id := root.ID()
	root.Children(func() {
		for i, name := range names {
			c.Key(i)
			underlined(c, parts.Tab(i), name, i == shown, id, nil)
		}
	})
	if kit.Changed(root, sound.Select) {
		kit.Log("UnderlineTabs", "picked", "index", *active)
	}
	return root
}
