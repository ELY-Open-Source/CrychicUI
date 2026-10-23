package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Share is one named fraction.
type Share struct {
	Name  string
	Share float64
}

// Traffic shows sources with shares.
func Traffic(c *ui.Context, sources []Share) ui.Element {
	if len(sources) == 0 {
		panic("dashboard: Traffic needs sources")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, source := range sources {
			if source.Name == "" || source.Share < 0 || source.Share > 1 || !finite(source.Share) {
				panic("dashboard: Traffic has an invalid source")
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, source.Name).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", source.Share*100)).TextColor(t.TextMuted)
				})
				progress(c, source.Share, t.Accent)
			})
		}
	})
	return root
}
