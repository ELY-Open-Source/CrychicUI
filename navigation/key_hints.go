package navigation

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// KeyHint is a shortcut's keys and what it does.
type KeyHint struct {
	Keys  []string
	Label string
}

// KeyHints is a footer of shortcut hints.
func KeyHints(c *ui.Context, label string, hints []KeyHint) ui.Element {
	if label == "" || len(hints) == 0 {
		panic(fmt.Sprintf("navigation: KeyHints %q needs hints", label))
	}
	for _, h := range hints {
		if len(h.Keys) == 0 || h.Label == "" {
			panic(fmt.Sprintf("navigation: KeyHints %q hints need keys and a label", label))
		}
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(4)).Padding(t.Space(2), t.Space(3)).AlignItems(ui.Center).
		Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Label(label)
	return root.Children(func() {
		for i, h := range hints {
			c.Key(i)
			hint := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
			kit.Arrive(hint, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*50*time.Millisecond).Children(func() {
				for j, k := range h.Keys {
					c.Key(j)
					typography.Kbd(c, k, false)
				}
				kit.Muted(c, h.Label).Padding(0, 0, 0, t.Space(1))
			})
		}
	})
}
