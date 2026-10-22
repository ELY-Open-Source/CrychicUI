package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// InlineMessage is tone icon and text fading in.
func InlineMessage(c *ui.Context, tone Tone, text string) ui.Element {
	if text == "" {
		panic("feedback: InlineMessage needs text")
	}
	t := c.Theme()
	col, glyph := look(t, tone)
	root := ui.Row(c).Gap(t.Space(1.5)).AlignSelf(ui.Start).Role(ui.RoleStatus).Label(text)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(col)
		ui.Text(c, text).TextColor(t.Text.Mix(col, 0.35))
	})
}
