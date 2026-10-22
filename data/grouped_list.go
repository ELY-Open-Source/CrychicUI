package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Group is a titled group of items.
type Group struct {
	Title string
	Items []string
}

// GroupedListOptions tune a GroupedList.
type GroupedListOptions struct {
	// Height of the list; zero gives the default.
	Height float32
}

// GroupedList pins group titles; chosen is group, item.
func GroupedList(c *ui.Context, label string, groups []Group, chosen *[2]int, o GroupedListOptions) ui.Element {
	if label == "" || len(groups) == 0 || chosen == nil {
		panic(fmt.Sprintf("data: GroupedList %q needs groups and a choice", label))
	}
	var flat [][2]int
	titles := map[string]bool{}
	for g, grp := range groups {
		if grp.Title == "" || titles[grp.Title] {
			panic(fmt.Sprintf("data: GroupedList %q groups need unique titles", label))
		}
		titles[grp.Title] = true
		flat = append(flat, [2]int{g, -1})
		for i, it := range grp.Items {
			if it == "" {
				panic(fmt.Sprintf("data: GroupedList %q items need text", label))
			}
			flat = append(flat, [2]int{g, i})
		}
	}
	sel := -1
	for r, at := range flat {
		if at == *chosen && at[1] >= 0 {
			sel = r
		}
	}
	if sel < 0 && *chosen != [2]int{-1, -1} {
		panic(fmt.Sprintf("data: GroupedList %q choice %v is not an item", label, *chosen))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Height(height(c, o.Height))
	list := ui.Local(root, "list", func() ui.ListState { return ui.ListState{} })
	list.Selected, list.RoundRows = &sel, true
	list.Key = func(r int) any { return flat[r] }
	list.Header = func(r int) bool { return flat[r][1] < 0 }
	root.Children(func() {
		l := ui.List(c, list, len(flat), func(r int) {
			at := flat[r]
			if at[1] < 0 {
				ui.Row(c).FillWidth().Padding(t.Space(1.5), t.Space(3)).Background(t.Surface).Children(func() {
					kit.Muted(c, groups[at[0]].Title).Bold().FontSize(t.Rem(0.75))
				})
				return
			}
			ui.Row(c).FillWidth().Padding(t.Space(2), t.Space(3)).Children(func() {
				ui.Text(c, groups[at[0]].Items[at[1]]).SingleLine()
			})
		}).Label(label).Grow(1)
		if l.Changed() && sel >= 0 && flat[sel][1] >= 0 {
			*chosen = flat[sel]
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("GroupedList", "chose", "group", flat[sel][0], "index", flat[sel][1])
		}
	})
	return root
}
