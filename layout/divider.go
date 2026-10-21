package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DividerOptions tune a Divider.
type DividerOptions struct {
	// Label sits on the line's center when set.
	Label string
}

// Divider is a hairline, labeled or not, growing from center.
func Divider(c *ui.Context, o DividerOptions) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(3)).Children(func() {
		hair := func(key string) {
			line := ui.Box(c).Key(key).Grow(1).Basis(0).Height(1).Background(t.Border)
			line.Scale(grow(line), 1)
		}
		hair("start")
		if o.Label != "" {
			kit.Muted(c, o.Label).FontSize(t.Rem(0.85))
			hair("end")
		}
	})
}
