package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Folder is a folder and unread count; Icon optional.
type Folder struct {
	ID, Name, Icon string
	Count          int
}

// Folders selects mail folders with counts; cue Select.
func Folders(c *ui.Context, folders []Folder, selected *string) ui.Element {
	if len(folders) == 0 || selected == nil {
		panic("mail: Folders needs folders and selection")
	}
	ids := make([]string, len(folders))
	for i, folder := range folders {
		if folder.Name == "" || folder.Count < 0 {
			panic(fmt.Sprintf("mail: Folders has an invalid folder %q", folder.ID))
		}
		ids[i] = folder.ID
	}
	unique("Folders", ids)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, folder := range folders {
			c.Key(folder.ID)
			on := *selected == folder.ID
			mark := folder.Icon
			if mark == "" {
				mark = "folder"
			}
			r := row(c, folder.Name, on)
			r.TextColor(r.AnimateColor("ink", t.Text.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, r, i).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, mark, 16)
					name := ui.Text(c, folder.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
					if on {
						name.Bold()
					}
					if folder.Count > 0 {
						ink := t.TextMuted
						if on {
							ink = t.Accent
						}
						kit.Chip(c, fmt.Sprint(folder.Count), ink)
					}
				})
			})
			if r.Clicked() && !on {
				*selected = folder.ID
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Folders", "selected", "index", i)
			}
		}
	})
	return root
}
