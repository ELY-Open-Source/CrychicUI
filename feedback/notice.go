package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Notice is a centered message that scales in.
func Notice(c *ui.Context, glyph, title, body string) ui.Element {
	if title == "" {
		panic("feedback: Notice needs a title")
	}
	t := c.Theme()
	svg := icon.Must(glyph)
	root := ui.Column(c).AlignSelf(ui.Center).MaxWidth(t.Space(100)).AlignItems(ui.Center).Gap(t.Space(2)).
		Padding(t.Space(6), t.Space(8)).Radius(t.Radius + 2).Background(t.Surface).Role(ui.RoleStatus).Label(title)
	kit.Arrive(root, ui.Pose{Scale: 0.94, Fade: 1}, 0)
	return root.Children(func() {
		ui.Icon(c, svg).Size(16, 16).TextColor(t.TextMuted)
		ui.Text(c, title).Bold().TextAlign(ui.Center)
		if body != "" {
			kit.Muted(c, body).TextAlign(ui.Center)
		}
	})
}
