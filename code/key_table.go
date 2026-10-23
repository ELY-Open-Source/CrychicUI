// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"strings"
	"time"
)

type Shortcut struct {
	Keys, Action string
}

// KeyGroup is a named group of shortcuts.

type KeyGroup struct {
	Name      string
	Shortcuts []Shortcut
}

// KeyTable is keyboard shortcuts in key caps by group.

func KeyTable(c *ui.Context, label string, groups []KeyGroup) ui.Element {
	if label == "" || len(groups) == 0 {
		panic("code: KeyTable needs a label and groups")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(3)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for gi, g := range groups {
			if g.Name == "" || len(g.Shortcuts) == 0 {
				panic("code: KeyTable needs a name and shortcuts in every group")
			}
			c.Key(g.Name)
			col := ui.Column(c).Gap(t.Space(1))
			kit.Arrive(col, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(gi*60)*time.Millisecond)
			col.Children(func() {
				kit.Muted(c, g.Name).FontSize(t.Rem(0.75)).Bold()
				for _, s := range g.Shortcuts {
					if s.Keys == "" || s.Action == "" {
						panic("code: KeyTable needs keys and an action in every shortcut")
					}
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
					row.Children(func() {
						keys := ui.Row(c).Gap(t.Space(1)).Width(t.Space(40)).Shrink(0)
						keys.Children(func() {
							for ki, k := range strings.Split(s.Keys, "+") {
								c.Key(ki)
								ui.Text(c, k).Font("monospace").FontSize(t.Rem(0.75)).SingleLine().
									Padding(t.Space(0.5), t.Space(1.5)).Background(t.Surface).Border(1, t.Border).Radius(t.Space(1))
							}
						})
						ui.Text(c, s.Action).FontSize(t.Rem(0.85)).SingleLine()
					})
				}
			})
		}
	})
	return root
}

// LogLevel is a console level.
