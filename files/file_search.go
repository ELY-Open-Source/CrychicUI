package files

import (
	"path/filepath"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileSearchOptions are the search labels.
type FileSearchOptions struct {
	Placeholder, Empty, All string
}

// FileSearch filters names and types; cue Select.
func FileSearch(c *ui.Context, label string, list []FileItem, query, kind *string, selected *int, o FileSearchOptions) ui.Element {
	validFiles("FileSearch", label, list)
	if query == nil || kind == nil || selected == nil || *selected < -1 || *selected >= len(list) {
		panic("files: FileSearch needs bound search values")
	}
	def(&o.Placeholder, "Search files")
	def(&o.Empty, "No files match")
	def(&o.All, "All")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field := input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		if field.Changed() {
			root.SetChanged()
		}
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			filters := []string{o.All}
			for _, item := range list {
				ext := strings.ToLower(filepath.Ext(item.Name))
				if ext != "" && !slices.Contains(filters, ext) {
					filters = append(filters, ext)
				}
			}
			for _, filter := range filters {
				active := *kind == filter || filter == o.All && *kind == ""
				chip := kit.Press(ui.ButtonBase(c).Label(filter).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(active))
				chip.Background(chip.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(active)), ui.Snappy))
				chip.Children(func() { ui.Text(c, filter).FontSize(t.Rem(0.75)).SingleLine() })
				if chip.Clicked() {
					next := filter
					if filter == o.All {
						next = ""
					}
					if next != *kind {
						*kind = next
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("FileSearch", "filter", "kind", *kind)
					}
				}
			}
		})
		q := strings.ToLower(*query)
		hits := 0
		for i, item := range list {
			if *kind != "" && !strings.EqualFold(filepath.Ext(item.Name), *kind) {
				continue
			}
			if q != "" && !strings.Contains(strings.ToLower(item.Name), q) {
				continue
			}
			hits++
			c.Key(i)
			row := fileButton(c, item.Name, i == *selected)
			row.Children(func() { ui.Text(c, item.Name).SingleLine() })
			if row.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("FileSearch", "selected", "index", i)
			}
		}
		if hits == 0 {
			kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
		}
	})
	return root
}
