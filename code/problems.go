package code

import (
	"fmt"
	"strconv"

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

// ProblemsOptions word a line, "%d" taking the number.
type ProblemsOptions struct {
	Line string
}

// Problems lists errors and warnings by file under their tallies.
func Problems(c *ui.Context, label string, problems []Problem, selected *int, o ProblemsOptions) ui.Element {
	if label == "" || len(problems) == 0 || selected == nil || *selected < -1 || *selected >= len(problems) {
		panic(fmt.Sprintf("code: Problems %q needs a label, problems and a choice or -1", label))
	}
	var count [severityCount]int
	for _, p := range problems {
		if p.Severity >= severityCount || p.File == "" || p.Line < 1 || p.Message == "" {
			panic(fmt.Sprintf("code: Problems %q needs a severity, a file, a line and a message", label))
		}
		count[p.Severity]++
	}
	def(&o.Line, "Ln %d")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Padding(0, 0, t.Space(1), 0).Role(ui.RoleList).Label(label)
	var order []string
	byFile := map[string][]int{}
	for i, p := range problems {
		if _, ok := byFile[p.File]; !ok {
			order = append(order, p.File)
		}
		byFile[p.File] = append(byFile[p.File], i)
	}
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
			tally := func(glyph string, ink ui.Color, n int) {
				ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(ink).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(14, 14)
					kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.75)).Bold().Shrink(0)
				})
			}
			tally("circle-x", t.Danger, count[ProbError])
			tally("triangle-alert", t.Warning, count[ProbWarning])
		})
		ui.Column(c).FillWidth().Padding(t.Space(1)).FocusGroup(ui.Vertical).Children(func() {
			k := 0
			for _, file := range order {
				c.Key(file)
				ui.Column(c).FillWidth().Children(func() { k = problemGroup(c, root, file, byFile[file], problems, selected, o.Line, k) })
			}
		})
	})
}

// problemGroup builds a file's rows; it returns the stagger.
func problemGroup(c *ui.Context, root ui.Element, file string, rows []int, problems []Problem, selected *int, line string, k int) int {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2), t.Space(0.5), t.Space(2)).Children(func() {
		ui.Icon(c, icon.Must(glyphFor(file))).Size(14, 14).TextColor(t.TextMuted)
		ui.Text(c, file).Font("monospace").FontSize(t.Rem(0.8)).Bold().Shrink(0)
		kit.Numeric(c, strconv.Itoa(len(rows))).FontSize(t.Rem(0.7)).Bold().TextColor(t.TextMuted).
			Padding(0, t.Space(1.5)).Radius(999).Background(t.Surface)
	})
	seen := map[string]int{}
	for _, i := range rows {
		p := problems[i]
		c.Key(ident(seen, fmt.Sprintf("%d %d %q", p.Severity, p.Line, p.Message)))
		problemRow(c, root, i, k, p, selected, line)
		k++
	}
	return k
}

// problemRow builds one selectable problem row.
func problemRow(c *ui.Context, root ui.Element, i, k int, p Problem, selected *int, line string) {
	t := c.Theme()
	on := i == *selected
	glyph, ink := "triangle-alert", t.Warning
	if p.Severity == ProbError {
		glyph, ink = "circle-x", t.Danger
	}
	b := kit.Press(ui.ButtonBase(c).Label(p.Message).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2), t.Space(1.25), t.Space(6)).
		Radius(t.Radius).Checked(on))
	kit.Stagger(c, root, b, k)
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
	b.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0).TextColor(ink)
			ui.Text(c, p.Message).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, fmt.Sprintf(line, p.Line)).Font("monospace").FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
		})
	})
	if b.Clicked() && !on {
		*selected = i
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("Problems", "selected", "index", i)
	}
}
