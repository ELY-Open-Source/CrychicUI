package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DocCardOptions are the document labels.
type DocCardOptions struct{ LastEdit string }

// DocCard shows a document with icon and last edit.
func DocCard(c *ui.Context, title, lastEdit string, o DocCardOptions) ui.Element {
	if title == "" || lastEdit == "" {
		panic("project: DocCard needs title and edit time")
	}
	def(&o.LastEdit, "Edited %s")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("file-text")).Size(t.Space(8), t.Space(8)).TextColor(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, title).Bold()
			kit.Muted(c, fmt.Sprintf(o.LastEdit, lastEdit)).FontSize(t.Rem(0.75))
		})
	})
	return root
}
