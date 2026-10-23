package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HandoffOptions word the event, "%s" taking the name.
type HandoffOptions struct {
	Joined string
}

// Handoff is an "agent joined" event with an avatar.
func Handoff(c *ui.Context, name string, image *ui.Bitmap, at time.Time, o HandoffOptions) ui.Element {
	if name == "" || at.IsZero() {
		panic(fmt.Sprintf("chat: Handoff %q needs a name and a time", name))
	}
	def(&o.Joined, "%s joined the chat")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Justify(ui.Center).Role(ui.RoleStatus)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		display.Avatar(c, name, image)
		ui.Column(c).Children(func() {
			ui.Text(c, sprintf(o.Joined, name)).FontSize(t.Rem(0.9))
			kit.Muted(c, when(c, at)).FontSize(t.Rem(0.75))
		})
	})
}
