package mail

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Folder is one mail folder and unread count.
type Folder struct {
	ID, Name string
	Count    int
}

// Folders selects mail folders with counts; cue Select.
func Folders(c *ui.Context, folders []Folder, selected *string) ui.Element {
	if len(folders) == 0 || selected == nil {
		panic("mail: Folders needs folders and selection")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, folder := range folders {
			if folder.ID == "" || folder.Name == "" || folder.Count < 0 {
				panic("mail: invalid folder")
			}
			c.Key(folder.ID)
			row := button(c, folder.Name).Row().FillWidth().AlignItems(ui.Center).Checked(*selected == folder.ID)
			row.Background(row.AnimateColor("selected", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*selected == folder.ID)), ui.Snappy))
			row.Children(func() {
				ui.Text(c, folder.Name).Grow(1)
				kit.Numeric(c, fmt.Sprint(folder.Count)).TextColor(t.TextMuted)
			})
			if row.Clicked() && *selected != folder.ID {
				*selected = folder.ID
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Folders", "selected")
			}
		}
	})
	return root
}
