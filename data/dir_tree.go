package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Dir is a folder, its count and subfolders.
type Dir struct {
	Name     string
	Count    int
	Children []Dir
}

// DirTreeOptions tune a DirTree.
type DirTreeOptions struct {
	// Expand names an arrow; it takes the name.
	Expand string
}

// DirTree is a folder tree; paths join with "/".
func DirTree(c *ui.Context, label string, dirs []Dir, open *ui.Selection[string], chosen *string, o DirTreeOptions) ui.Element {
	if label == "" || len(dirs) == 0 || open == nil || chosen == nil {
		panic(fmt.Sprintf("data: DirTree %q needs folders, open folders and a choice", label))
	}
	def(&o.Expand, "Expand %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Role(ui.RoleTree).Label(label)
	shown := *chosen
	var level func(ds []Dir, at string, depth int)
	level = func(ds []Dir, at string, depth int) {
		var names []string
		for _, d := range ds {
			if d.Count < 0 {
				panic(fmt.Sprintf("data: DirTree %q counts must not be negative", label))
			}
			names = append(names, d.Name)
		}
		named("DirTree", label, names)
		for _, d := range ds {
			path := at + d.Name
			c.Key(d.Name)
			ui.Column(c).FillWidth().Children(func() {
				branch(c, root, "DirTree", fmt.Sprintf(o.Expand, d.Name), path, depth, len(d.Children) > 0, open, func() {
					b := ui.ButtonBase(c).Label(d.Name).Grow(1).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).
						Radius(t.Radius).Checked(path == shown)
					bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
					if path == shown {
						bg = t.Accent.Alpha(0.12)
					}
					b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
						glyph := "folder"
						if open.Has(path) {
							glyph = "folder-open"
						}
						ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(t.Accent).ContentTransition(ui.ContentTransition{})
						ui.Text(c, d.Name).Grow(1).SingleLine()
						kit.Numeric(c, fmt.Sprint(d.Count)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
					})
					if b.Clicked() {
						*chosen = path
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("DirTree", "chose", "depth", depth)
					}
				})
				kids(c, open.Has(path), func() { level(d.Children, path+"/", depth+1) })
			})
		}
	}
	root.Children(func() { level(dirs, "", 0) })
	return root
}
