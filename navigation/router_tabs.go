package navigation

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RouterTabs are tabs over r's pages, which slide.
func RouterTabs(c *ui.Context, label string, r *ui.Router, pages []Page, page func(rt *ui.Route)) ui.Element {
	at := routed("RouterTabs", label, r, pages)
	if page == nil {
		panic(fmt.Sprintf("navigation: RouterTabs %q needs a page builder", label))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	last := ui.Local(root, "last", func() int { return at })
	run := ui.Local(root, "run", func() uint64 { return 0 })
	dir := ui.Local(root, "dir", func() float32 { return 0 })
	if at != *last {
		*dir = 1
		if (at < *last) != (c.Direction() == ui.RightToLeft) {
			*dir = -1
		}
		*last = at
		*run++
	}
	return root.Children(func() {
		active := at
		parts := ui.TabsBase(c, &active, len(pages))
		list := parts.List.Label(label).Gap(t.Space(1))
		id := list.ID()
		list.Children(func() {
			for i, p := range pages {
				c.Key(i)
				underlined(c, parts.Tab(i), p.Name, i == at, id, nil)
			}
		})
		if kit.Changed(list, sound.Select) {
			r.Push(pages[active].Path)
			root.SetChanged()
			kit.Log("RouterTabs", "pushed", "index", active)
		}
		ui.Box(c).FillWidth().Clip().Children(func() {
			body := ui.Column(c).FillWidth()
			// Replays start from the shown value; jump first.
			x := body.Keyframes("slide", []ui.Keyframe{{At: 0, Value: 0}, {At: time.Millisecond, Value: *dir * t.Space(8)},
				{At: 250 * time.Millisecond, Value: 0, Ease: ui.EaseOut}}, ui.KeyframeOptions{Run: *run})
			body.Translate(x, 0).Children(func() { r.View(c, page) })
		})
	})
}
