package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ArriveList slides new items in; the rest make room.
func ArriveList(c *ui.Context, items []string) ui.Element {
	if len(items) > 0 {
		unique("ArriveList", items)
	}
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(2))
	seen := ui.Local(root, "seen", func() map[string]bool { return set(items) })
	fresh := 0
	for _, s := range items {
		if !(*seen)[s] {
			fresh++
		}
	}
	*seen = set(items)
	if fresh > 0 {
		sound.Play(sound.Notify)
		kit.Log("ArriveList", "arrived", "count", fresh)
	}
	arrive := ui.ElementTransition{Animation: &ui.Snappy, Enter: &ui.Motion{Y: -t.Space(4)}, Exit: &ui.Motion{Collapse: true}}
	return root.Children(func() {
		for _, s := range items {
			c.Key(s)
			ui.Box(c).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
				Transition(arrive).Children(func() { ui.Text(c, s) })
		}
	})
}
