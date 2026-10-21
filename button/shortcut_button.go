package button

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ShortcutButtonOptions tune a ShortcutButton.
type ShortcutButtonOptions struct {
	// Disabled blocks click and shortcut activation.
	Disabled bool
}

// ShortcutButton shows keys and fires on click or shortcut.
func ShortcutButton(c *ui.Context, label string, mods ui.Modifiers, key ui.Key, keys []string, o ShortcutButtonOptions) ui.Element {
	if label == "" || len(keys) == 0 {
		panic("button: ShortcutButton needs a label and keys")
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface).Disabled(o.Disabled)
	root.Children(func() {
		ui.Text(c, label).SingleLine()
		kit.Muted(c, strings.Join(keys, "+")).Padding(t.Space(0.5), t.Space(1)).Radius(999).
			Background(t.Accent.Alpha(0.12)).TextColor(t.Accent)
	})
	if !o.Disabled && !root.IsDisabled() && (root.Clicked() || c.Shortcut(mods, key)) {
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("ShortcutButton", "activated")
	}
	return root
}
