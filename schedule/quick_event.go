package schedule

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuickEventOptions are the quick event labels.
type QuickEventOptions struct{ Placeholder, Create string }

// QuickEvent enters a one-line event; cue Pop.
func QuickEvent(c *ui.Context, title *string, created *bool, o QuickEventOptions) (root, create ui.Element) {
	if title == nil || created == nil {
		panic("schedule: QuickEvent needs title and state")
	}
	def(&o.Placeholder, "Event title")
	def(&o.Create, "Create")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	pending := ui.Local(root, "pending", func() bool { return false })
	root.Children(func() {
		in := ui.TextInput(c, title).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
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
			Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*title) == ""))
		create.Children(func() { ui.Text(c, o.Create).SingleLine() })
	})
	if *pending || create.Clicked() {
		*pending = false
		*created = true
		create.SetChanged()
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("QuickEvent", "created", "runes", utf8.RuneCountInString(*title))
	}
	return root, create
}
