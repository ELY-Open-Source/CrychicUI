package navigation

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Sitemap is a nested page list fading in.
func Sitemap(c *ui.Context, label string, list []Node, picked *string) ui.Element {
	if label == "" || len(list) == 0 || picked == nil {
		panic(fmt.Sprintf("navigation: Sitemap %q needs pages and a pick", label))
	}
	nodes("Sitemap", list)
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Label(label)
	shown, n := *picked, 0
	var level func(list []Node, parent string, depth int)
	level = func(list []Node, parent string, depth int) {
		for _, node := range list {
			path := node.Label
			if parent != "" {
				path = parent + "/" + node.Label
			}
			c.Key(path)
			b := rowBase(c, node.Label).Checked(path == shown).Padding(t.Space(1), t.Space(3), t.Space(1), t.Space(3)+float32(depth)*t.Space(5))
			fg := t.Text
			if path == shown {
				fg = t.Accent
			}
			kit.Arrive(b, ui.Pose{X: -t.Space(2), Fade: 1}, time.Duration(n)*40*time.Millisecond)
			n++
			b.TextColor(b.AnimateColor("fg", fg, ui.Snappy)).Children(func() {
				text := ui.Text(c, node.Label).SingleLine().FillWidth()
				if depth == 0 {
					text.Bold()
				}
			})
			if kit.Clicked(b, sound.Tap) && path != *picked {
				*picked = path
				root.SetChanged()
				kit.Log("Sitemap", "picked", "depth", depth)
			}
			level(node.Children, path, depth+1)
		}
	}
	return root.Children(func() { level(list, "", 0) })
}
