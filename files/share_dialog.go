package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SharePerson is one person, their email and access role.
type SharePerson struct {
	Name, Email, Role string
}

// ShareDialogOptions word the dialog.
type ShareDialogOptions struct {
	Title, Copy, Anyone string
}

// ShareDialog shows people with their roles and a copyable link.
func ShareDialog(c *ui.Context, label string, people []SharePerson, link string, o ShareDialogOptions) (root, copy ui.Element) {
	if label == "" || len(people) == 0 || link == "" {
		panic(fmt.Sprintf("files: ShareDialog %q needs people and a link", label))
	}
	seen := map[string]bool{}
	for _, p := range people {
		if p.Name == "" || p.Role == "" || seen[p.Name+"\x00"+p.Email] {
			panic(fmt.Sprintf("files: ShareDialog %q people need unique names and roles", label))
		}
		seen[p.Name+"\x00"+p.Email] = true
	}
	def(&o.Title, "People with access")
	def(&o.Copy, "Copy link")
	def(&o.Anyone, "Anyone with the link can view")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleDialog).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Scale: 0.98, Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "share-2", label, nil)
		ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).Children(func() {
			kit.Muted(c, o.Title).FontSize(t.Rem(0.7)).Bold().Padding(0, 0, t.Space(1), 0)
			for i, p := range people {
				c.Key(fmt.Sprintf("%q %q", p.Name, p.Email))
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1), 0)
				kit.Stagger(c, root, row, i).Children(func() {
					display.Avatar(c, p.Name, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.7)).Shrink(0)
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, p.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().SingleLine()
						if p.Email != "" {
							kit.Muted(c, p.Email).FontSize(t.Rem(0.75)).FillWidth().SingleLine()
						}
					})
					kit.Chip(c, p.Role, t.TextMuted)
				})
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2.5), t.Space(3)).Background(t.Surface).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(t.Space(6), t.Space(6)).Center().Radius(999).Background(t.Success.Alpha(0.14)).TextColor(t.Success).
					Children(func() { ui.Icon(c, icon.Must("globe")).Size(14, 14) })
				ui.Text(c, o.Anyone).FontSize(t.Rem(0.8))
			})
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(2.5)).
				Background(t.Background).Border(1, t.Border).Radius(t.Radius).Children(func() {
				ui.Icon(c, icon.Must("link")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
				ui.Text(c, link).Font("monospace").FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).SingleLine()
				copy = button.CopyButton(c, o.Copy, link)
			})
		})
	})
	if copy.Clicked() {
		kit.Log("ShareDialog", "copied", "people", len(people))
	}
	return root, copy
}
