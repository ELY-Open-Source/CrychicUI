package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Severity is a problem's level.
type Severity uint8

// Problem severities.
const (
	ProbError Severity = iota
	ProbWarning
	severityCount
)

// Problem is one error or warning.
type Problem struct {
	Severity Severity
	File     string
	Line     int
	Message  string
}

// Problems lists errors and warnings by file.
func Problems(c *ui.Context, label string, problems []Problem, selected *int) ui.Element {
	if label == "" || len(problems) == 0 || selected == nil || *selected < -1 || *selected >= len(problems) {
		panic(fmt.Sprintf("code: Problems %q needs a label, problems and a choice or -1", label))
	}
	for _, p := range problems {
		if p.Severity >= severityCount || p.File == "" || p.Line < 1 || p.Message == "" {
			panic(fmt.Sprintf("code: Problems %q needs a severity, a file, a line and a message", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	var order []string
	byFile := map[string][]int{}
	for i, p := range problems {
		if _, ok := byFile[p.File]; !ok {
			order = append(order, p.File)
		}
		byFile[p.File] = append(byFile[p.File], i)
	}
	return root.Children(func() {
		for _, file := range order {
			c.Key("f." + file)
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), 0, t.Space(1)).Children(func() {
				ui.Icon(c, icon.Must("file-code")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, file).Bold().FontSize(t.Rem(0.85)).SingleLine()
			})
			for _, i := range byFile[file] {
				c.Key(i)
				problemRow(c, root, i, problems[i], selected)
			}
		}
	})
}

// problemRow builds one selectable problem row.
func problemRow(c *ui.Context, root ui.Element, i int, p Problem, selected *int) {
	t := c.Theme()
	on := i == *selected
	glyph, ink := "triangle-alert", t.Warning
	if p.Severity == ProbError {
		glyph, ink = "circle-x", t.Danger
	}
	b := kit.Press(ui.ButtonBase(c).Label(p.Message).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Margin(0, 0, 0, t.Space(4)).Checked(on))
	bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
	if on {
		bg = t.Accent.Alpha(0.1)
	}
	b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).TextColor(ink)
		ui.Text(c, p.Message).Grow(1).SingleLine()
		kit.Muted(c, fmt.Sprintf(":%d", p.Line)).Font("monospace").FontSize(t.Rem(0.8))
	})
	if b.Clicked() && !on {
		*selected = i
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("Problems", "selected", "index", i)
	}
}
