package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SystemMessage is a centered muted event line.
func SystemMessage(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("chat: SystemMessage needs text")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Justify(ui.Center).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		kit.Muted(c, text).FontSize(t.Rem(0.8)).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Surface)
	})
}
