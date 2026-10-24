package files

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// flash is how long a rename glows.
const flash = 900 * time.Millisecond

// RenameRow edits a name in place and glows; cue Success.
func RenameRow(c *ui.Context, name *string) ui.Element {
	if name == nil || strings.TrimSpace(*name) == "" {
		panic("files: RenameRow needs a bound nonempty name")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius).Border(1, t.Border)
	last := ui.Local(root, "last", func() string { return *name })
	saved := ui.Local(root, "saved", func() time.Time { return time.Time{} })
	glow := c.Now().Sub(*saved) < flash
	if glow {
		c.After(flash - c.Now().Sub(*saved))
	}
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.Success.Alpha(0.14), unit(glow)), ui.Smooth))
	root.Children(func() {
		fileTile(c, FileItem{Name: *name}, t.Space(8))
		edit := ui.EditableText(c, name).Grow(1).MinWidth(0)
		if edit.Changed() {
			if strings.TrimSpace(*name) == "" {
				*name = *last
				kit.Log("RenameRow", "empty rejected")
			} else {
				*last = *name
				*saved = c.Now()
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("RenameRow", "renamed", "runes", utf8.RuneCountInString(*name))
			}
		} else {
			*last = *name
		}
		mark := ui.Icon(c, icon.Must("pencil")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
		mark.Opacity(mark.AnimateTo("hover", 0.4+0.6*unit(root.Hovered()), ui.Snappy))
	})
	return root
}
