package display

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TagOptions tune a Tag.
type TagOptions struct {
	// Icon names an optional leading Lucide icon.
	Icon string
}

// Tag is a rounded label fading in.
func Tag(c *ui.Context, text string, o TagOptions) ui.Element {
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
		Background(t.Surface).Border(1, t.Border).Children(func() {
		if o.Icon != "" {
			ui.Icon(c, icon.Must(o.Icon)).Size(16, 16)
		}
		ui.Text(c, text).FontSize(t.Rem(0.85))
	})
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0)
}
