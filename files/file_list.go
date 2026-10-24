package files

import (
	"cmp"
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileSort selects the file list ordering.
type FileSort uint8

// File list orders.
const (
	SortName FileSort = iota
	SortSize
	SortDate
)

// FileListOptions name the columns and word dates.
type FileListOptions struct {
	Name, Size, Date string
	FormatDate       func(t time.Time) string
}

var rowSlide = ui.ElementTransition{Animation: &ui.Snappy, Position: true}

// FileList shows file rows that sort by column; cue Select.
func FileList(c *ui.Context, label string, list []FileItem, selected *int, sort *FileSort, o FileListOptions) ui.Element {
	validFiles("FileList", label, list)
	unique("FileList", label, list)
	if selected == nil || sort == nil || *selected < -1 || *selected >= len(list) || *sort > SortDate {
		panic(fmt.Sprintf("files: FileList %q needs a selection or -1 and a known sort", label))
	}
	def(&o.Name, "Name")
	def(&o.Size, "Size")
	def(&o.Date, "Modified")
	if o.FormatDate == nil {
		o.FormatDate = shortDate
	}
	t := c.Theme()
	ordered := make([]int, len(list))
	for i := range ordered {
		ordered[i] = i
	}
	slices.SortStableFunc(ordered, func(a, b int) int {
		fa, fb := list[a], list[b]
		if fa.Dir != fb.Dir {
			return cmp.Compare(unit(fb.Dir), unit(fa.Dir))
		}
		switch *sort {
		case SortSize:
			return cmp.Compare(fa.Size, fb.Size)
		case SortDate:
			return fa.Modified.Compare(fb.Modified)
		}
		return cmp.Compare(fa.Name, fb.Name)
	})
	sizeW, dateW := t.Space(18), t.Space(26)
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(3)).Background(t.Surface).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			for _, h := range []struct {
				name  string
				sort  FileSort
				width float32
			}{{o.Name, SortName, 0}, {o.Size, SortSize, sizeW}, {o.Date, SortDate, dateW}} {
				c.Key(h.sort)
				on := *sort == h.sort
				b := kit.Press(ui.ButtonBase(c).Label(h.name).Padding(t.Space(0.5), t.Space(1)).Radius(t.Radius).Checked(on))
				if h.width > 0 {
					b.Width(h.width).Shrink(0)
				} else {
					b.Grow(1).MinWidth(0).Padding(t.Space(0.5), t.Space(1), t.Space(0.5), t.Space(10))
				}
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
				b.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, h.name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
						arrow := ui.Icon(c, icon.Must("arrow-down")).Size(12, 12)
						arrow.Opacity(arrow.AnimateTo("on", unit(on), ui.Snappy))
					})
				})
				if b.Clicked() && !on {
					*sort = h.sort
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("FileList", "sorted", "sort", int(*sort))
				}
			}
		})
		ui.Column(c).FillWidth().Padding(t.Space(1)).FocusGroup(ui.Vertical).Children(func() {
			for k, i := range ordered {
				f := list[i]
				c.Key(fileKey(f))
				on := i == *selected
				row := kit.Press(ui.ButtonBase(c).Label(f.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
					Radius(t.Radius).Checked(on)).Transition(rowSlide)
				kit.Stagger(c, root, row, k)
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						fileTile(c, f, t.Space(7))
						name := ui.Text(c, f.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						if on {
							name.Bold()
						}
						size := ui.Box(c).Width(sizeW).Shrink(0)
						size.Children(func() {
							if f.Dir {
								kit.Muted(c, "—").FontSize(t.Rem(0.8))
							} else {
								sizeText(c, f.Size)
							}
						})
						kit.Muted(c, o.FormatDate(f.Modified)).FontSize(t.Rem(0.8)).Width(dateW).Shrink(0).SingleLine()
					})
				})
				if row.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("FileList", "selected", "index", i)
				}
			}
		})
	})
	return root
}
