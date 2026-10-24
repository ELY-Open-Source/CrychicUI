package files

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NewFolderOptions word the new folder row.
type NewFolderOptions struct {
	Label, Placeholder, Create, Invalid string
}

// NewFolder names and submits a new folder; cue Pop.
func NewFolder(c *ui.Context, name *string, o NewFolderOptions) (root, create ui.Element) {
	if name == nil {
		panic("files: NewFolder needs a bound name")
	}
	def(&o.Label, "New folder")
	def(&o.Placeholder, "Folder name")
	def(&o.Create, "Create")
	def(&o.Invalid, `Names cannot contain / or \`)
	t := c.Theme()
	root = kit.Card(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2.5), t.Space(2.5), t.Space(2.5), t.Space(3)).Role(ui.RoleGroup).Label(o.Label)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	pending := ui.Local(root, "pending", func() bool { return false })
	valid := func() bool {
		s := strings.TrimSpace(*name)
		return s != "" && !strings.ContainsAny(s, `/\`)
	}
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(t.Warning.Alpha(0.14)).TextColor(t.Warning).
				Children(func() { ui.Icon(c, icon.Must("folder-plus")).Size(16, 16) })
			in := ui.TextInput(c, name).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || k.Mods != 0 {
					return false
				}
				*pending = true
				return true
			})
			if in.Changed() {
				root.SetChanged()
			}
			create = kit.Press(ui.ButtonBase(c).Label(o.Create).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius)).Disabled(!valid())
			create.Background(create.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(create.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() { ui.Text(c, o.Create).Bold().Shrink(0) })
		})
		if strings.ContainsAny(*name, `/\`) {
			ui.Row(c).Key("invalid").Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(0, 0, 0, t.Space(10.5)).Transition(kit.Enter).Children(func() {
				ui.Icon(c, icon.Must("circle-x")).Size(12, 12).TextColor(t.Danger)
				ui.Text(c, o.Invalid).FontSize(t.Rem(0.75)).TextColor(t.Danger)
			})
		}
	})
	submit := *pending || create.Clicked()
	*pending = false
	if submit && valid() {
		create.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("NewFolder", "created", "runes", utf8.RuneCountInString(*name))
	}
	return root, create
}
