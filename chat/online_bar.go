package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OnlineBar shows avatars of who is online; new ones pop.
func OnlineBar(c *ui.Context, label string, names []string) ui.Element {
	if label == "" || len(names) == 0 {
		panic(fmt.Sprintf("chat: OnlineBar %q needs a label and names", label))
	}
	seen := map[string]bool{}
	for _, n := range names {
		if n == "" || seen[n] {
			panic(fmt.Sprintf("chat: OnlineBar %q names must be unique and not empty", label))
		}
		seen[n] = true
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(label)
	known := ui.Local(root, "known", func() map[string]bool { return map[string]bool{} })
	return root.Children(func() {
		for _, n := range names {
			c.Key(n)
			_, was := (*known)[n]
			row := ui.Box(c)
			if *ui.Local(row, "fresh", func() bool { return !was }) {
				kit.Arrive(row, ui.Pose{Scale: 0.5, Fade: 1}, 0)
			}
			row.Children(func() {
				display.AvatarStatus(c, n, nil, display.PresenceOnline)
			})
		}
		for _, n := range names {
			(*known)[n] = true
		}
		for n := range *known {
			if !seen[n] {
				delete(*known, n)
			}
		}
	})
}
