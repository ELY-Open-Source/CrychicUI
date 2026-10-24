package files

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileSearchOptions word the search.
type FileSearchOptions struct {
	Placeholder, Empty, All string
}

// FileSearch filters files by name and type; cue Select.
func FileSearch(c *ui.Context, label string, list []FileItem, query, kind *string, selected *int, o FileSearchOptions) ui.Element {
	validFiles("FileSearch", label, list)
	unique("FileSearch", label, list)
	if query == nil || kind == nil || selected == nil || *selected < -1 || *selected >= len(list) {
		panic(fmt.Sprintf("files: FileSearch %q needs bound search values", label))
	}
	def(&o.Placeholder, "Search files")
	def(&o.Empty, "No files match")
	def(&o.All, "All")
	t := c.Theme()
	counts := map[string]int{}
	var kinds []string
	for _, f := range list {
		ext := strings.ToLower(filepath.Ext(f.Name))
		if ext == "" {
			continue
		}
		if counts[ext] == 0 {
			kinds = append(kinds, ext)
		}
		counts[ext]++
	}
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			if input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder}).FillWidth().Changed() {
				root.SetChanged()
			}
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
				for _, k := range append([]string{""}, kinds...) {
					c.Key("kind" + k)
					name, n := k, counts[k]
					if k == "" {
						name, n = o.All, len(list)
					}
					on := *kind == k
					chip := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(0.75), t.Space(2.5)).Radius(999).Checked(on))
					chip.Background(chip.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(chip.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
					chip.TextColor(chip.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
					chip.Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, name).FontSize(t.Rem(0.75)).Bold().Shrink(0).MaxWidth(t.Space(36)).SingleLine()
							kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.7)).Shrink(0).Opacity(0.7)
						})
					})
					if chip.Clicked() && !on {
						*kind = k
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("FileSearch", "filter", "kind", k)
					}
				}
			})
		})
		q := strings.TrimSpace(*query)
		hits := 0
		ui.Column(c).FillWidth().Padding(t.Space(1)).FocusGroup(ui.Vertical).Role(ui.RoleList).Children(func() {
			for i, f := range list {
				if *kind != "" && !strings.EqualFold(filepath.Ext(f.Name), *kind) {
					continue
				}
				lo, hi := kit.IndexFold(f.Name, q)
				if q != "" && lo < 0 {
					continue
				}
				c.Key(fileKey(f))
				on := i == *selected
				row := kit.Press(ui.ButtonBase(c).Label(f.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
					Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, row, hits)
				hits++
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						fileTile(c, f, t.Space(7))
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.RichText(c, kit.Marked(t, f.Name, lo, hi)...).FontSize(t.Rem(0.85)).FillWidth().SingleLine()
							if f.Path != "" {
								kit.Muted(c, f.Path).FontSize(t.Rem(0.7)).FillWidth().SingleLine()
							}
						})
						sizeText(c, f.Size).FontSize(t.Rem(0.75))
					})
				})
				if row.Clicked() && !on {
					*selected = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("FileSearch", "selected", "index", i)
				}
			}
			if hits == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("search-x")).Size(20, 20).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
		})
	})
	return root
}
