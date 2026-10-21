package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Section is a muted heading, growing divider and body.
func Section(c *ui.Context, title string, content func()) ui.Element {
	t := c.Theme()
	return ui.Column(c).Gap(t.Space(2)).Children(func() {
		kit.Muted(c, title).Bold().FontSize(t.Rem(0.85))
		line := ui.Box(c).FillWidth().Height(1).Background(t.Border).TransformOrigin(0, 0.5)
		line.Scale(grow(line), 1)
		ui.Column(c).Gap(t.Space(2)).Children(content)
	})
}
