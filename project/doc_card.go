package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DocCardOptions word the edit line; Author and Preview are optional.
type DocCardOptions struct{ LastEdit, By, Author, Preview, Icon string }

// DocCard shows a document with icon and last edit.
func DocCard(c *ui.Context, title, lastEdit string, o DocCardOptions) ui.Element {
	if title == "" || lastEdit == "" {
		panic(fmt.Sprintf("project: DocCard %q needs a title and edit time", title))
	}
	def(&o.LastEdit, "Edited %s")
	def(&o.By, "Edited %s by %s")
	def(&o.Icon, "file-text")
	t := c.Theme()
	meta := fmt.Sprintf(o.LastEdit, lastEdit)
	if o.Author != "" {
		meta = fmt.Sprintf(o.By, lastEdit, o.Author)
	}
	root := kit.Lift(card(c).Gap(t.Space(2.5)).Padding(t.Space(3)))
	hint := root.AnimateTo("hint", unit(root.Hovered()), ui.Snappy)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			tile(c, o.Icon, t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.75)).Children(func() {
				ui.Text(c, title).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					if o.Author != "" {
						display.Avatar(c, o.Author, nil).Size(t.Space(5), t.Space(5)).FontSize(t.Rem(0.5))
					}
					kit.Muted(c, meta).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
				})
			})
			glyph(c, "arrow-up-right", 16).TextColor(t.Accent).Opacity(hint).Translate(t.Space(1)*(1-hint), 0)
		})
		if o.Preview != "" {
			kit.Muted(c, o.Preview).FontSize(t.Rem(0.75)).Wrap().FillWidth().Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
				Background(t.SurfaceHover.Alpha(0.6))
		}
	})
	return root
}
