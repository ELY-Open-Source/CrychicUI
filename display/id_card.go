package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// IDCard pairs an avatar and name with field rows.
func IDCard(c *ui.Context, name string, image *ui.Bitmap, fields []KeyValuePair) ui.Element {
	if len(fields) == 0 {
		panic("display: IDCard needs fields")
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4))
	kit.Arrive(root, ui.Pose{Y: 8, Fade: 1}, 0).Children(func() {
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			avatarCircle(c, name, image, t.Space(10))
			ui.Text(c, name).Bold().FontSize(t.Rem(1.25))
		})
		ui.Column(c).Gap(t.Space(2)).Children(func() {
			for _, field := range fields {
				KeyValue(c, field.Key, field.Value)
			}
		})
	})
	return root
}
