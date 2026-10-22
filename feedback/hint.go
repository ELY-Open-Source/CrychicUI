package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Hint is muted help text with an info icon.
func Hint(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("feedback: Hint needs text")
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Start).TextColor(t.TextMuted)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		ui.Icon(c, icon.Must("info")).Size(16, 16).Margin(t.Space(0.25), 0, 0, 0)
		ui.Text(c, text).FontSize(t.Rem(0.85))
	})
}
