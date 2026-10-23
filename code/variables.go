package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Variable is a name, type, value and children.
type Variable struct {
	Name, Type, Value string
	Children          []Variable
}

// Variables is a tree of names, types and values.
func Variables(c *ui.Context, label string, vars []Variable) ui.Element {
	if label == "" {
		panic("code: Variables needs a label")
	}
	checkVariables(label, vars)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Role(ui.RoleTree).Label(label)
	root.Children(func() {
		variableRows(c, root, vars, 0, "")
	})
	return root
}

// checkVariables panics on a variable without a name.
func checkVariables(label string, vars []Variable) {
	for _, v := range vars {
		if v.Name == "" {
			panic(fmt.Sprintf("code: Variables %q variables need names", label))
		}
		checkVariables(label, v.Children)
	}
}

// variableRows builds rows of vars at depth under path.
func variableRows(c *ui.Context, root ui.Element, vars []Variable, depth int, path string) {
	t := c.Theme()
	for i, v := range vars {
		key := fmt.Sprintf("%s.%d", path, i)
		c.Key(key)
		if len(v.Children) > 0 {
			btn := kit.Press(ui.ButtonBase(c).Label(v.Name).Row().FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(0, t.Space(1)).Margin(0, 0, 0, float32(depth)*t.Space(4)))
			open := ui.Local(btn, "open", func() bool { return depth == 0 })
			btn.Background(btn.AnimateColor("bg", t.SurfaceHover.Alpha(unit(btn.Hovered())), ui.Snappy)).Children(func() {
				ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
				ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
				variableText(c, v)
			})
			if btn.Clicked() {
				*open = !*open
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Variables", "open", *open)
			}
			if *open {
				variableRows(c, root, v.Children, depth+1, key)
			}
		} else {
			ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(0, t.Space(1)).Margin(0, 0, 0, float32(depth)*t.Space(4)).Children(func() {
				ui.Box(c).Size(16, 16)
				variableText(c, v)
			})
		}
	}
}

// variableText builds a variable's name, type and value.
func variableText(c *ui.Context, v Variable) {
	t := c.Theme()
	ui.Text(c, v.Name).Font("monospace").FontSize(t.Rem(0.85)).SingleLine()
	if v.Type != "" {
		kit.Muted(c, v.Type).FontSize(t.Rem(0.8)).SingleLine()
	}
	ui.Text(c, v.Value).Font("monospace").FontSize(t.Rem(0.85)).Bold().Grow(1).SingleLine()
}
