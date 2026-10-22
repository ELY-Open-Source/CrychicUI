package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/motion"
	"github.com/egoist/mygo/ui"
)

// Stat is a label, a whole value and a unit.
type Stat struct {
	Label string
	Value int64
	Unit  string
}

// StatList shows label and value rows; values roll.
func StatList(c *ui.Context, label string, stats []Stat) ui.Element {
	if label == "" || len(stats) == 0 {
		panic(fmt.Sprintf("data: StatList %q needs stats", label))
	}
	var names []string
	for _, s := range stats {
		names = append(names, s.Label)
	}
	unique("StatList", label, names)
	t := c.Theme()
	root := frame(c, label).Dividers(1, t.Border)
	return root.Children(func() {
		for _, s := range stats {
			c.Key(s.Label)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).Children(func() {
				kit.Muted(c, s.Label).Grow(1)
				motion.RollingDigits(c, s.Value).FontSize(t.Rem(1))
				if s.Unit != "" {
					kit.Muted(c, s.Unit)
				}
			})
		}
	})
}
