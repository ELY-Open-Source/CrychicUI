package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PermissionRole is a file access level.
type PermissionRole uint8

// Access levels.
const (
	Viewer PermissionRole = iota
	Editor
	Owner
)

// PermissionsOptions name the roles and say what each may do.
type PermissionsOptions struct {
	Viewer, Editor, Owner             string
	ViewerNote, EditorNote, OwnerNote string
}

// Permissions picks owner, editor or viewer; cue Select.
func Permissions(c *ui.Context, label string, role *PermissionRole, o PermissionsOptions) ui.Element {
	if label == "" || role == nil || *role > Owner {
		panic(fmt.Sprintf("files: Permissions %q needs a label and a known role", label))
	}
	def(&o.Viewer, "Viewer")
	def(&o.Editor, "Editor")
	def(&o.Owner, "Owner")
	def(&o.ViewerNote, "Can view and comment")
	def(&o.EditorNote, "Can edit and share")
	def(&o.OwnerNote, "Full control, including deleting")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).FocusGroup(ui.Vertical).Role(ui.RoleRadioGroup).Label(label)
	root.Children(func() {
		for i, it := range []struct {
			role              PermissionRole
			name, note, glyph string
		}{{Owner, o.Owner, o.OwnerNote, "crown"}, {Editor, o.Editor, o.EditorNote, "pencil"}, {Viewer, o.Viewer, o.ViewerNote, "eye"}} {
			c.Key(it.role)
			on := *role == it.role
			row := kit.Press(ui.ButtonBase(c).Label(it.name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(2.5), t.Space(3)).
				Radius(t.Radius).Checked(on).Role(ui.RoleRadio))
			kit.Stagger(c, root, row, i)
			row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())).Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
			row.Border(1, row.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			row.Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					ink := t.TextMuted.Mix(t.Accent, unit(on))
					ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink).
						Children(func() { ui.Icon(c, icon.Must(it.glyph)).Size(16, 16) })
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, it.name).Bold().FontSize(t.Rem(0.9))
						kit.Muted(c, it.note).FontSize(t.Rem(0.8)).Wrap()
					})
					ring := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Radius(999)
					ring.Border(t.Space(0.5), ring.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
					dot := ring.AnimateTo("dot", unit(on), ui.Bouncy)
					ring.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
						d := r.W * 0.5 * dot
						p.Fill(ui.Rect{X: r.X + (r.W-d)/2, Y: r.Y + (r.H-d)/2, W: d, H: d}, t.Accent, 999)
					}))
				})
			})
			if row.Clicked() && !on {
				*role = it.role
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Permissions", "selected", "role", int(*role))
			}
		}
	})
	return root
}
