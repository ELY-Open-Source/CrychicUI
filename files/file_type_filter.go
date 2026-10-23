package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TypeCount is one file type and its count.
type TypeCount struct {
	Kind  string
	Count int
}

// FileTypeFilter toggles type chips with counts; cue Toggle.
func FileTypeFilter(c *ui.Context, label string, types []TypeCount, selected *int) ui.Element {
	if label == "" || len(types) == 0 || selected == nil || *selected < -1 || *selected >= len(types) {
		panic("files: FileTypeFilter needs types and selection")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(1.5)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, item := range types {
			if item.Kind == "" || item.Count < 0 {
				panic("files: FileTypeFilter has an invalid type")
			}
			c.Key(i)
			name := fmt.Sprintf("%s %d", item.Kind, item.Count)
			b := kit.Press(ui.ButtonBase(c).Label(name).Row().Gap(t.Space(1)).Padding(t.Space(1), t.Space(2)).
				Radius(999).Border(1, t.Border).Checked(i == *selected))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(i == *selected)), ui.Snappy))
			b.Children(func() {
				ui.Text(c, item.Kind).FontSize(t.Rem(0.8))
				kit.Numeric(c, fmt.Sprint(item.Count)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			})
			if b.Clicked() {
				if *selected == i {
					*selected = -1
				} else {
					*selected = i
				}
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("FileTypeFilter", "selected", "index", *selected)
			}
		}
	})
	return root
}
