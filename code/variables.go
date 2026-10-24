package code

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Variable is a name, type, value and children.
type Variable struct {
	Name, Type, Value string
	Children          []Variable
}

// Variables is a titled tree of names, values and types.
func Variables(c *ui.Context, label string, vars []Variable) ui.Element {
	if label == "" || len(vars) == 0 {
		panic(fmt.Sprintf("code: Variables %q needs a label and variables", label))
	}
	checkVariables(label, vars)
	t := c.Theme()
	root := kit.Card(c).FillWidth().Padding(t.Space(1), t.Space(1), t.Space(2), t.Space(1)).Role(ui.RoleTree).Label(label)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Children(func() {
			ui.Icon(c, icon.Must("braces")).Size(14, 14).TextColor(t.TextMuted)
			ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Numeric(c, fmt.Sprint(len(vars))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
		})
		variableRows(c, root, vars, 0)
	})
	return root
}

// checkVariables panics on a variable without a name.
func checkVariables(label string, vars []Variable) {
	seen := map[string]bool{}
	for _, v := range vars {
		if v.Name == "" || seen[v.Name] {
			panic(fmt.Sprintf("code: Variables %q variables need unique names", label))
		}
		seen[v.Name] = true
		checkVariables(label, v.Children)
	}
}

// variableRows builds rows of vars at depth; children fold.
func variableRows(c *ui.Context, root ui.Element, vars []Variable, depth int) {
	t := c.Theme()
	for i, v := range vars {
		c.Key(v.Name)
		if len(v.Children) == 0 {
			row := ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(0.75), t.Space(2)).Radius(t.Radius).
				Role(ui.RoleTreeItem).Label(v.Name)
			kit.Arrive(row, ui.Pose{X: -t.Space(1), Fade: 1}, time.Duration(i)*30*time.Millisecond).Children(func() {
				ui.Box(c).Size(14, 14).Shrink(0)
				variableText(c, v)
			})
			continue
		}
		col := ui.Column(c).FillWidth()
		open := ui.Local(col, "open", func() bool { return depth == 0 })
		var btn ui.Element
		col.Children(func() {
			btn = kit.Press(ui.ButtonBase(c).Label(v.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(0.75), t.Space(2)).
				Radius(t.Radius).Role(ui.RoleTreeItem).Expanded(*open))
			btn.Background(btn.AnimateColor("bg", t.SurfaceHover.Alpha(unit(btn.Hovered())), ui.Snappy)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ch := ui.Icon(c, icon.Must("chevron-right")).Size(14, 14).Shrink(0).TextColor(t.TextMuted)
					ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
					variableText(c, v)
				})
			})
			layout.Collapse(c, *open, func() {
				ui.Column(c).FillWidth().Padding(0, 0, 0, t.Space(3.5)).Children(func() {
					ui.Column(c).FillWidth().Padding(0, 0, 0, t.Space(1)).BorderWidth(0, 0, 0, 1).BorderColor(t.Border).Children(func() {
						variableRows(c, root, v.Children, depth+1)
					})
				})
			})
		})
		if btn.Clicked() {
			*open = !*open
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("Variables", "open", *open, "depth", depth)
		}
	}
}

// variableText builds a variable's name, value and type.
func variableText(c *ui.Context, v Variable) {
	t := c.Theme()
	ui.Text(c, v.Name).Font("monospace").FontSize(t.Rem(0.8)).TextColor(t.Accent.Mix(t.Text, 0.35)).Shrink(0)
	kit.Muted(c, "=").Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
	ui.Text(c, v.Value).Font("monospace").FontSize(t.Rem(0.8)).TextColor(valueInk(t, v.Value)).Grow(1).MinWidth(0).SingleLine()
	if v.Type != "" {
		kit.Muted(c, v.Type).Font("monospace").FontSize(t.Rem(0.7)).Shrink(0)
	}
}
