package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecentFiles shows recently opened rows; cue Tap.
func RecentFiles(c *ui.Context, label string, list []FileItem, opened *int) ui.Element {
	validFiles("RecentFiles", label, list)
	if opened == nil || *opened < -1 || *opened >= len(list) {
		panic("files: RecentFiles needs a valid open index")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, item := range list {
			c.Key(i)
			row := fileButton(c, item.Name, i == *opened)
			row.Children(func() {
				ui.Icon(c, icon.Must(itemGlyph(item))).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, item.Name).SingleLine().Grow(1)
				ui.Text(c, date(item.Modified)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
			})
			if row.Clicked() {
				*opened = i
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("RecentFiles", "opened", "index", i)
			}
		}
	})
	return root
}
