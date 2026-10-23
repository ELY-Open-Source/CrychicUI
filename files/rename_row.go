package files

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RenameRow edits a filename in place; cue Success.
func RenameRow(c *ui.Context, name *string) ui.Element {
	if name == nil || strings.TrimSpace(*name) == "" {
		panic("files: RenameRow needs a bound nonempty name")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).
		Radius(t.Radius).Border(1, t.Border)
	last := ui.Local(root, "last", func() string { return *name })
	root.Children(func() {
		ui.Icon(c, icon.Must(glyph(*name))).Size(16, 16).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{})
		edit := ui.EditableText(c, name).Grow(1)
		if edit.Changed() {
			if strings.TrimSpace(*name) == "" {
				*name = *last
				kit.Log("RenameRow", "empty rejected")
				return
			}
			*last = *name
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("RenameRow", "renamed", "runes", utf8.RuneCountInString(*name))
		} else {
			*last = *name
		}
	})
	return root
}
