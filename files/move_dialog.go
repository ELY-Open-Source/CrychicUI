package files

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MoveDialogOptions are the move labels.
type MoveDialogOptions struct {
	Title, Move string
}

// MoveDialog chooses a folder target; cue Drop.
func MoveDialog(c *ui.Context, label string, folders []string, target *int, o MoveDialogOptions) (root, move ui.Element) {
	if label == "" || len(folders) == 0 || target == nil || *target < -1 || *target >= len(folders) {
		panic("files: MoveDialog needs folders and a valid target")
	}
	def(&o.Title, "Choose a destination")
	def(&o.Move, "Move")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(64)).Gap(t.Space(2)).Padding(t.Space(3)).
		Radius(t.Radius).Border(1, t.Border).Background(t.Surface).Role(ui.RoleDialog).Label(label)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		for i, folder := range folders {
			if folder == "" {
				panic("files: MoveDialog folder must not be empty")
			}
			c.Key(i)
			depth := strings.Count(folder, "/")
			name := folder
			if at := strings.LastIndex(folder, "/"); at >= 0 {
				name = folder[at+1:]
			}
			row := fileButton(c, folder, i == *target).Padding(t.Space(1), t.Space(float32(2+2*depth)))
			row.Children(func() {
				ui.Icon(c, icon.Must("folder")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, name).SingleLine()
			})
			if row.Clicked() && *target != i {
				*target = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("MoveDialog", "target", "index", i)
			}
		}
		move = kit.Press(ui.ButtonBase(c).Label(o.Move).AlignSelf(ui.End).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText).Disabled(*target < 0))
		move.Children(func() { ui.Text(c, o.Move).SingleLine() })
	})
	if move.Clicked() && *target >= 0 {
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("MoveDialog", "moved", "target", *target)
	}
	return root, move
}
