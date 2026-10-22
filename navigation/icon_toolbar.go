package navigation

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IconToolbar is groups of icon buttons with tooltips.
func IconToolbar(c *ui.Context, label string, groups [][]NavItem, used *[2]int) ui.Element {
	if label == "" || len(groups) == 0 || used == nil {
		panic(fmt.Sprintf("navigation: IconToolbar %q needs groups and a used pointer", label))
	}
	for g, tools := range groups {
		if len(tools) == 0 {
			panic(fmt.Sprintf("navigation: IconToolbar %q group %d is empty", label, g))
		}
		for _, tool := range tools {
			if tool.Label == "" {
				panic(fmt.Sprintf("navigation: IconToolbar %q tools need labels", label))
			}
			icon.Must(tool.Icon)
		}
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Padding(t.Space(1)).AlignItems(ui.Center).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal).Label(label)
	return root.Children(func() {
		for g, tools := range groups {
			if g > 0 {
				c.Key(fmt.Sprint("gap", g))
				ui.Box(c).Width(1).Height(t.Space(5)).Background(t.Border)
			}
			for i, tool := range tools {
				c.Key(fmt.Sprint(g, "/", i))
				b := kit.Press(ui.ButtonBase(c).Label(tool.Label).Tooltip(tool.Label).Size(t.Space(8), t.Space(8)).Center().Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy)).Children(func() {
					ui.Icon(c, icon.Must(tool.Icon)).Size(16, 16)
				})
				if kit.Clicked(b, sound.Tap) {
					*used = [2]int{g, i}
					root.SetChanged()
					kit.Log("IconToolbar", "used", "group", g, "index", i)
				}
			}
		}
	})
}
