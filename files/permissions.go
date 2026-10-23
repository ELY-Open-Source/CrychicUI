package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PermissionRole is a file access level.
type PermissionRole uint8

const (
	Viewer PermissionRole = iota
	Editor
	Owner
)

// PermissionsOptions are the role labels.
type PermissionsOptions struct {
	Viewer, Editor, Owner string
}

// Permissions selects owner, editor, or viewer; cue Select.
func Permissions(c *ui.Context, label string, role *PermissionRole, o PermissionsOptions) ui.Element {
	if label == "" || role == nil || *role > Owner {
		panic("files: Permissions needs a label and valid role")
	}
	def(&o.Viewer, "Viewer")
	def(&o.Editor, "Editor")
	def(&o.Owner, "Owner")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for _, item := range []struct {
			role PermissionRole
			name string
		}{{Owner, o.Owner}, {Editor, o.Editor}, {Viewer, o.Viewer}} {
			row := fileButton(c, item.name, *role == item.role)
			row.Children(func() {
				ui.Icon(c, icon.Must("users")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, item.name).Grow(1)
				if *role == item.role {
					ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Accent)
				}
			})
			if row.Clicked() && *role != item.role {
				*role = item.role
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Permissions", "selected", "role", *role)
			}
		}
	})
	return root
}
