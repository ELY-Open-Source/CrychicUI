package files

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NewFolderOptions are the new folder labels.
type NewFolderOptions struct {
	Label, Placeholder, Create string
}

// NewFolder names and submits a folder; cue Pop.
func NewFolder(c *ui.Context, name *string, o NewFolderOptions) (root, create ui.Element) {
	if name == nil {
		panic("files: NewFolder needs a bound name")
	}
	def(&o.Label, "New folder")
	def(&o.Placeholder, "Folder name")
	def(&o.Create, "Create")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(o.Label)
	pending := ui.Local(root, "pending", func() bool { return false })
	root.Children(func() {
		in := ui.TextInput(c, name).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
		in.OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
				return false
			}
			*pending = true
			return true
		})
		if in.Changed() {
			root.SetChanged()
		}
		create = kit.Press(ui.ButtonBase(c).Label(o.Create).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*name) == ""))
		create.Children(func() { ui.Text(c, o.Create).SingleLine() })
	})
	if *pending || create.Clicked() {
		*pending = false
		create.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("NewFolder", "created", "runes", utf8.RuneCountInString(*name))
	}
	return root, create
}
