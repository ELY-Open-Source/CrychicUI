package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TrashOptions are the trash labels.
type TrashOptions struct {
	Restore, Empty string
}

// Trash shows deleted items with restore; cue Success.
func Trash(c *ui.Context, label string, items []string, restored *int, o TrashOptions) ui.Element {
	if label == "" || restored == nil || *restored < -1 || *restored >= len(items) {
		panic("files: Trash needs a label and valid restore index")
	}
	def(&o.Restore, "Restore")
	def(&o.Empty, "Trash is empty")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		if len(items) == 0 {
			kit.Muted(c, o.Empty).Padding(t.Space(3))
		}
		for i, name := range items {
			if name == "" {
				panic("files: Trash item must not be empty")
			}
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).
				Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(row, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
			row.Children(func() {
				ui.Icon(c, icon.Must("trash-2")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, name).Grow(1).SingleLine()
				b := kit.Press(ui.ButtonBase(c).Label(o.Restore).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
				b.Children(func() { ui.Text(c, o.Restore).FontSize(t.Rem(0.8)).SingleLine() })
				if b.Clicked() {
					*restored = i
					root.SetChanged()
					sound.Play(sound.Success)
					kit.Log("Trash", "restored", "index", i)
				}
			})
		}
	})
	return root
}
