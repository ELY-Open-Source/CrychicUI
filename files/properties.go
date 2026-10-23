package files

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PropertiesOptions are the property labels.
type PropertiesOptions struct {
	Size, Type, Modified, Created, Path string
}

// Properties shows file size, type, dates, and path.
func Properties(c *ui.Context, item FileItem, path, kind string, created time.Time, o PropertiesOptions) ui.Element {
	validFile("Properties", item)
	if path == "" || kind == "" {
		panic("files: Properties needs a path and type")
	}
	def(&o.Size, "Size")
	def(&o.Type, "Type")
	def(&o.Modified, "Modified")
	def(&o.Created, "Created")
	def(&o.Path, "Path")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, row := range []struct{ name, value string }{
			{o.Size, sizeLabel(item.Size)},
			{o.Type, kind},
			{o.Modified, date(item.Modified)},
			{o.Created, date(created)},
			{o.Path, path},
		} {
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				kit.Muted(c, row.name).FontSize(t.Rem(0.8)).Width(t.Space(18))
				ui.Text(c, row.value).FontSize(t.Rem(0.85)).Wrap()
			})
		}
	})
	return root
}
