package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LinkPreview is a title, host and image card; fades in.
func LinkPreview(c *ui.Context, title, host string, image *ui.Bitmap) ui.Element {
	if title == "" || host == "" {
		panic(fmt.Sprintf("chat: LinkPreview %q needs a title and a host", title))
	}
	t := c.Theme()
	root := kit.Card(c).AlignSelf(ui.Start).Width(t.Space(70)).Clip().Role(ui.RoleGroup).Label(title)
	kit.Arrive(root, ui.Pose{Fade: 1}, 0)
	return root.Children(func() {
		if image != nil {
			ui.Image(c, image).FillWidth().AspectRatio(1.9).Fit(ui.Cover)
		}
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2), t.Space(3)).Children(func() {
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("link")).Size(16, 16).TextColor(t.TextMuted)
				kit.Muted(c, host).FontSize(t.Rem(0.8)).SingleLine()
			})
			ui.Text(c, title).Bold().Wrap()
		})
	})
}
