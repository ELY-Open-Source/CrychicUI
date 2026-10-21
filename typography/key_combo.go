package typography

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// KeyCombo presses its caps in sequence on a match.
func KeyCombo(c *ui.Context, labels []string, matched bool) ui.Element {
	if len(labels) == 0 {
		panic("typography: KeyCombo needs labels")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1))
	run := pressRun(root, matched)
	return root.Children(func() {
		for i, label := range labels {
			if i > 0 {
				kit.Muted(c, "+")
			}
			c.Key(i)
			keyCap(c, label, run, time.Duration(i)*60*time.Millisecond)
		}
	})
}
