package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PropertiesOptions name the fields; Bytes words the exact size.
type PropertiesOptions struct {
	Size, Type, Modified, Created, Path, Bytes string
	FormatDate                                 func(t time.Time) string
}

// Properties shows a file's size, type, dates and path.
func Properties(c *ui.Context, item FileItem, path, kind string, created time.Time, o PropertiesOptions) ui.Element {
	validFile("Properties", item)
	if path == "" || kind == "" || created.IsZero() {
		panic(fmt.Sprintf("files: Properties %q needs a path, a type and a creation time", item.Name))
	}
	def(&o.Size, "Size")
	def(&o.Type, "Type")
	def(&o.Modified, "Modified")
	def(&o.Created, "Created")
	def(&o.Path, "Path")
	def(&o.Bytes, "%d bytes")
	if o.FormatDate == nil {
		o.FormatDate = func(t time.Time) string { return t.Format("Jan 2, 2006 at 15:04") }
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(item.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			fileTile(c, item, t.Space(11))
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, item.Name).Bold().FillWidth().SingleLine()
				kit.Muted(c, kind).FontSize(t.Rem(0.8))
			})
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			field := func(i int, name string, value func()) {
				c.Key(i)
				ui.Row(c).FillWidth().Gap(t.Space(3)).Children(func() {
					kit.Muted(c, name).FontSize(t.Rem(0.8)).Width(t.Space(20)).Shrink(0)
					ui.Column(c).Grow(1).MinWidth(0).Children(value)
				})
			}
			field(0, o.Size, func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					sizeText(c, item.Size).FontSize(t.Rem(0.85)).TextColor(t.Text)
					kit.Muted(c, fmt.Sprintf(o.Bytes, item.Size)).FontSize(t.Rem(0.75))
				})
			})
			field(1, o.Type, func() { ui.Text(c, kind).FontSize(t.Rem(0.85)) })
			field(2, o.Modified, func() { ui.Text(c, o.FormatDate(item.Modified)).FontSize(t.Rem(0.85)) })
			field(3, o.Created, func() { ui.Text(c, o.FormatDate(created)).FontSize(t.Rem(0.85)) })
			field(4, o.Path, func() { ui.Text(c, path).Font("monospace").FontSize(t.Rem(0.8)).Wrap() })
		})
	})
	return root
}
