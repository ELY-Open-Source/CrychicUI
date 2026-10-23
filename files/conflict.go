package files

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConflictChoice is a conflict resolution.
type ConflictChoice uint8

const (
	ConflictNone ConflictChoice = iota
	KeepMine
	KeepTheirs
	KeepBoth
)

// ConflictOptions are the resolution labels.
type ConflictOptions struct{ Mine, Theirs, Both string }

// Conflict resolves a file conflict; cue Select.
func Conflict(c *ui.Context, label string, choice *ConflictChoice, o ConflictOptions) ui.Element {
	if label == "" || choice == nil || *choice > KeepBoth {
		panic("files: Conflict needs a label and valid choice")
	}
	def(&o.Mine, "Keep mine")
	def(&o.Theirs, "Keep theirs")
	def(&o.Both, "Keep both")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for _, item := range []struct {
			choice ConflictChoice
			name   string
		}{{KeepMine, o.Mine}, {KeepTheirs, o.Theirs}, {KeepBoth, o.Both}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Border(1, t.Border).Checked(*choice == item.choice))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*choice == item.choice)), ui.Snappy))
			b.Children(func() { ui.Text(c, item.name).SingleLine() })
			if b.Clicked() && *choice != item.choice {
				*choice = item.choice
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Conflict", "resolved", "choice", *choice)
			}
		}
	})
	return root
}
