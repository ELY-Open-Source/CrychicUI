package files

import (
	"cmp"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileSort selects the file list ordering.
type FileSort uint8

const (
	SortName FileSort = iota
	SortSize
	SortDate
)

// FileListOptions are the sort labels.
type FileListOptions struct {
	Name, Size, Date string
}

// FileList shows sortable file rows; cue Select.
func FileList(c *ui.Context, label string, list []FileItem, selected *int, sort *FileSort, o FileListOptions) ui.Element {
	validFiles("FileList", label, list)
	if selected == nil || sort == nil || *selected < -1 || *selected >= len(list) || *sort > SortDate {
		panic("files: FileList needs valid selection and sort")
	}
	def(&o.Name, "Name")
	def(&o.Size, "Size")
	def(&o.Date, "Date")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleGroup).Label(label)
	ordered := make([]int, len(list))
	for i := range ordered {
		ordered[i] = i
	}
	slices.SortStableFunc(ordered, func(a, b int) int {
		fa, fb := list[a], list[b]
		switch *sort {
		case SortSize:
			return cmp.Compare(fa.Size, fb.Size)
		case SortDate:
			return fa.Modified.Compare(fb.Modified)
		default:
			return cmp.Compare(fa.Name, fb.Name)
		}
	})
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			for _, h := range []struct {
				name string
				sort FileSort
			}{{o.Name, SortName}, {o.Size, SortSize}, {o.Date, SortDate}} {
				b := kit.Press(ui.ButtonBase(c).Label(h.name).Padding(t.Space(0.5), t.Space(1)).Radius(t.Radius))
				b.Children(func() { ui.Text(c, h.name).FontSize(t.Rem(0.75)).TextColor(t.TextMuted) })
				if b.Clicked() && *sort != h.sort {
					*sort = h.sort
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("FileList", "sorted", "sort", *sort)
				}
			}
		})
		for _, original := range ordered {
			f := list[original]
			c.Key(original)
			row := fileButton(c, f.Name, original == *selected)
			row.Children(func() {
				ui.Icon(c, icon.Must(itemGlyph(f))).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, f.Name).SingleLine().Grow(1)
				sizeText(c, f.Size)
				ui.Text(c, date(f.Modified)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
			})
			if row.Clicked() && *selected != original {
				*selected = original
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("FileList", "selected", "index", original)
			}
		}
	})
	return root
}
