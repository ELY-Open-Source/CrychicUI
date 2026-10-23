package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ZipEntry is one archived file and its compressed size.
type ZipEntry struct {
	Name       string
	Size       int64
	Compressed int64
}

// ZipListOptions are the archive labels.
type ZipListOptions struct{ Compressed string }

// ZipList expands archive entries; cue Toggle.
func ZipList(c *ui.Context, label string, entries []ZipEntry, opened *int, o ZipListOptions) ui.Element {
	if label == "" || len(entries) == 0 || opened == nil || *opened < -1 || *opened >= len(entries) {
		panic("files: ZipList needs entries and a valid open index")
	}
	def(&o.Compressed, "Compressed")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, entry := range entries {
			if entry.Name == "" || entry.Size < 0 || entry.Compressed < 0 {
				panic("files: ZipList has an invalid entry")
			}
			c.Key(i)
			row := fileButton(c, entry.Name, i == *opened)
			row.Children(func() {
				ui.Icon(c, icon.Must("file-archive")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, entry.Name).SingleLine().Grow(1)
				sizeText(c, entry.Size)
			})
			if row.Clicked() {
				if *opened == i {
					*opened = -1
				} else {
					*opened = i
				}
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("ZipList", "toggled", "index", *opened)
			}
			if *opened == i {
				ui.Row(c).FillWidth().Padding(t.Space(1), t.Space(4)).Children(func() {
					kit.Muted(c, o.Compressed).FontSize(t.Rem(0.8)).Grow(1)
					sizeText(c, entry.Compressed)
				}).Transition(kit.Enter)
			}
		}
	})
	return root
}
