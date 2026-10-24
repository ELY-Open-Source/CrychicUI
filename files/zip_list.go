package files

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ZipEntry is one archived file and its compressed size.
type ZipEntry struct {
	Name             string
	Size, Compressed int64
}

// ZipListOptions word an entry's details.
type ZipListOptions struct {
	Compressed, Original string
}

// ZipList lists archive entries; one opens its sizes; cue Toggle.
func ZipList(c *ui.Context, label string, entries []ZipEntry, opened *int, o ZipListOptions) ui.Element {
	if label == "" || len(entries) == 0 || opened == nil || *opened < -1 || *opened >= len(entries) {
		panic(fmt.Sprintf("files: ZipList %q needs entries and an open index or -1", label))
	}
	var total, packed int64
	seen := map[string]bool{}
	for _, e := range entries {
		if e.Name == "" || seen[e.Name] || e.Size < 0 || e.Compressed < 0 {
			panic(fmt.Sprintf("files: ZipList %q needs unique named entries with sizes", label))
		}
		if e.Size > math.MaxInt64-total || e.Compressed > math.MaxInt64-packed {
			panic(fmt.Sprintf("files: ZipList %q totals exceed int64", label))
		}
		seen[e.Name] = true
		total, packed = total+e.Size, packed+e.Compressed
	}
	def(&o.Compressed, "Compressed")
	def(&o.Original, "Original")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleList).Label(label)
	root.Children(func() {
		kit.CardHead(c, "archive", label, func() {
			sizeText(c, packed).FontSize(t.Rem(0.75))
			kit.Chip(c, ratio(total, packed), t.Success)
		})
		ui.Column(c).FillWidth().Padding(t.Space(1)).Children(func() {
			for i, e := range entries {
				c.Key(e.Name)
				on := i == *opened
				col := ui.Column(c).FillWidth()
				kit.Stagger(c, root, col, i).Children(func() {
					row := kit.Press(ui.ButtonBase(c).Label(e.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
						Radius(t.Radius).Expanded(on))
					row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
					row.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ch := ui.Icon(c, icon.Must("chevron-right")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
							ch.Rotate(90 * ch.AnimateTo("turn", unit(on), ui.Snappy))
							fileTile(c, FileItem{Name: e.Name}, t.Space(6))
							ui.Text(c, e.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							sizeText(c, e.Size)
						})
					})
					layout.Collapse(c, on, func() {
						ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(3), t.Space(2), t.Space(12)).Children(func() {
							detail := func(name string, bytes int64, share float64, ink ui.Color) {
								ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
									kit.Muted(c, name).FontSize(t.Rem(0.75)).Width(t.Space(20)).Shrink(0)
									ui.Box(c).Grow(1).Children(func() { progressBar(c, share, ink) })
									sizeText(c, bytes).FontSize(t.Rem(0.75)).Width(t.Space(14)).Shrink(0)
								})
							}
							share := 0.0
							if e.Size > 0 {
								share = min(1, float64(e.Compressed)/float64(e.Size))
							}
							detail(o.Original, e.Size, 1, t.Border)
							detail(o.Compressed, e.Compressed, share, t.Success)
						})
					})
					if row.Clicked() {
						*opened = i
						if on {
							*opened = -1
						}
						root.SetChanged()
						sound.Play(sound.Toggle)
						kit.Log("ZipList", "toggled", "index", *opened)
					}
				})
			}
		})
	})
	return root
}

// ratio words how much packed shrank, as "−62%".
func ratio(total, packed int64) string {
	if total <= 0 || packed >= total {
		return "0%"
	}
	return "−" + strconv.Itoa(int(100*(1-float64(packed)/float64(total))+0.5)) + "%"
}
