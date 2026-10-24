package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GroupData is one social group summary.
type GroupData struct {
	ID, Name, Privacy string
	Members           int
}

// GroupCardOptions are the group count labels.
type GroupCardOptions struct{ Members string }

// GroupCard shows group members and privacy.
func GroupCard(c *ui.Context, group GroupData, o GroupCardOptions) ui.Element {
	if group.ID == "" || group.Name == "" || group.Privacy == "" || group.Members < 0 {
		panic("social: invalid group")
	}
	def(&o.Members, "%d members")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, group.Name).Bold()
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Members, group.Members)).TextColor(t.TextMuted)
			kit.Muted(c, group.Privacy)
		})
	})
	return root
}
