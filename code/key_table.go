package code

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Shortcut is one keyboard shortcut, keys joined by "+".
type Shortcut struct {
	Keys, Action string
}

// KeyGroup is a named group of shortcuts.
type KeyGroup struct {
	Name      string
	Shortcuts []Shortcut
}

// KeyTable lists keyboard shortcuts as key caps by group.
func KeyTable(c *ui.Context, label string, groups []KeyGroup) ui.Element {
	if label == "" || len(groups) == 0 {
		panic(fmt.Sprintf("code: KeyTable %q needs a label and groups", label))
	}
	for _, g := range groups {
		if g.Name == "" || len(g.Shortcuts) == 0 {
			panic(fmt.Sprintf("code: KeyTable %q needs a name and shortcuts in every group", label))
		}
		for _, s := range g.Shortcuts {
			if s.Keys == "" || s.Action == "" || strings.Contains("+"+s.Keys+"+", "++") {
				panic(fmt.Sprintf("code: KeyTable %q needs keys and an action in every shortcut", label))
			}
		}
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		kit.CardHead(c, "keyboard", label, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(2), t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			for gi, g := range groups {
				c.Key(g.Name)
				col := ui.Column(c).FillWidth().Gap(t.Space(0.5))
				kit.Arrive(col, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(gi)*60*time.Millisecond).Children(func() {
					kit.Muted(c, g.Name).FontSize(t.Rem(0.7)).Bold().Padding(0, 0, t.Space(1), 0)
					for si, s := range g.Shortcuts {
						c.Key(si)
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), 0).
							BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(si < len(g.Shortcuts)-1))).Children(func() {
							ui.Text(c, s.Action).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap()
							ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).Children(func() {
								for ki, k := range strings.Split(s.Keys, "+") {
									c.Key(ki)
									ui.Text(c, k).Font("monospace").FontSize(t.Rem(0.75)).Bold().Shrink(0).MinWidth(t.Space(6)).TextAlign(ui.Center).
										Padding(t.Space(0.5), t.Space(1.5)).Background(t.Background).BorderWidth(1, 1, 2, 1).BorderColor(t.Border).
										Radius(t.Space(1))
								}
							})
						})
					}
				})
			}
		})
	})
	return root
}
