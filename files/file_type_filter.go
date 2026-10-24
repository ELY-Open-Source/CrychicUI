package files

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TypeCount is one file type, a sample name and count.
type TypeCount struct {
	Kind string
	// Sample names a file of the kind, choosing the icon.
	Sample string
	Count  int
}

// FileTypeFilter toggles one type chip with its count; cue Toggle.
func FileTypeFilter(c *ui.Context, label string, types []TypeCount, selected *int) ui.Element {
	if label == "" || len(types) == 0 || selected == nil || *selected < -1 || *selected >= len(types) {
		panic(fmt.Sprintf("files: FileTypeFilter %q needs types and a selection or -1", label))
	}
	seen := map[string]bool{}
	for _, it := range types {
		if it.Kind == "" || seen[it.Kind] || it.Count < 0 {
			panic(fmt.Sprintf("files: FileTypeFilter %q needs unique kinds and counts not below zero", label))
		}
		seen[it.Kind] = true
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).FocusGroup(ui.Horizontal).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, it := range types {
			c.Key(it.Kind)
			on := i == *selected
			f := FileItem{Name: it.Sample}
			ink := kindInk(t, f)
			b := kit.Press(ui.ButtonBase(c).Label(it.Kind+" "+strconv.Itoa(it.Count)).Tooltip(it.Kind).Padding(t.Space(1), t.Space(1.5), t.Space(1), t.Space(2)).
				Radius(999).Checked(on))
			kit.Stagger(c, root, b, i)
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(ink.Alpha(0.16), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(ink, unit(on)), ui.Snappy))
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Icon(c, icon.Must(itemGlyph(f))).Size(14, 14).Shrink(0).TextColor(ink)
					ui.Text(c, it.Kind).FontSize(t.Rem(0.8)).Bold().Shrink(0).MaxWidth(t.Space(40)).SingleLine()
					kit.Numeric(c, strconv.Itoa(it.Count)).FontSize(t.Rem(0.7)).Bold().Shrink(0).TextColor(t.TextMuted).
						Padding(0, t.Space(1.5)).Radius(999).Background(t.Surface)
				})
			})
			if b.Clicked() {
				*selected = i
				if on {
					*selected = -1
				}
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("FileTypeFilter", "selected", "index", *selected)
			}
		}
	})
	return root
}
