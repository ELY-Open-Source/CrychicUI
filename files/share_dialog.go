package files

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SharePerson is one person and their access role.
type SharePerson struct {
	Name, Role string
}

// ShareDialogOptions are the sharing labels.
type ShareDialogOptions struct {
	Title, Copy string
}

// ShareDialog shows people, roles, and a copyable link.
func ShareDialog(c *ui.Context, label string, people []SharePerson, link string, o ShareDialogOptions) (root, copy ui.Element) {
	if label == "" || len(people) == 0 || link == "" {
		panic("files: ShareDialog needs people and a link")
	}
	def(&o.Title, "People with access")
	def(&o.Copy, "Copy link")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(64)).Gap(t.Space(2)).Padding(t.Space(3)).
		Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Role(ui.RoleDialog).Label(label)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		for i, person := range people {
			if person.Name == "" || person.Role == "" {
				panic("files: ShareDialog people need names and roles")
			}
			c.Key(i)
			ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
				ui.Text(c, person.Name).Grow(1)
				kit.Muted(c, person.Role).FontSize(t.Rem(0.8))
			})
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, link).Font("monospace").FontSize(t.Rem(0.75)).SingleLine().Grow(1)
			copy = button.CopyButton(c, o.Copy, link)
		})
	})
	return root, copy
}
