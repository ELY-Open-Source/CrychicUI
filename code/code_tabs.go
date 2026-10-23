package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CodeFile is one tab of a CodeTabs view.
type CodeFile struct {
	Name, Lang, Code string
}

// CodeTabs shows files as tabs over a code block.
func CodeTabs(c *ui.Context, files []CodeFile, current *int) ui.Element {
	if len(files) == 0 || current == nil || *current < 0 || *current >= len(files) {
		panic(fmt.Sprintf("code: CodeTabs needs files and a tab, got %d", *current))
	}
	for _, f := range files {
		if f.Name == "" || f.Code == "" {
			panic("code: CodeTabs files need a name and code")
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Border(1, t.Border).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1), t.Space(1), 0, t.Space(1)).Background(t.Surface).FocusGroup(ui.Horizontal).Children(func() {
			for i, f := range files {
				c.Key(i)
				on := i == *current
				b := kit.Press(ui.ButtonBase(c).Label(f.Name).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(on))
				bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
				if on {
					bg = t.Background
				}
				b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
					ui.Text(c, f.Name).Font("monospace").FontSize(t.Rem(0.85)).SingleLine()
				})
				if b.Clicked() && !on {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("CodeTabs", "selected", "index", i)
				}
			}
		})
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		ui.Column(c).FillWidth().Padding(t.Space(2)).Background(t.Background).Children(func() {
			f := files[*current]
			codeLines(c, f.Lang, f.Code, nil)
		})
	})
	return root
}
