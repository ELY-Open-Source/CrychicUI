package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CallStack shows frames with the current one marked.
func CallStack(c *ui.Context, list []StackFrame, current *int) ui.Element {
	if len(list) == 0 || current == nil || *current < 0 || *current >= len(list) {
		panic(fmt.Sprintf("code: CallStack needs frames and a current one, got %d", *current))
	}
	frames("CallStack", list)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label("Call stack")
	return root.Children(func() {
		for i, f := range list {
			c.Key(i)
			on := i == *current
			b := kit.Press(ui.ButtonBase(c).Label(f.Func).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(on))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
				ui.Box(c).Width(t.Space(0.75)).Shrink(0).FillHeight().Radius(999).Background(t.Accent.Alpha(unit(on)))
				ui.Text(c, f.Func).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).SingleLine()
				kit.Muted(c, fmt.Sprintf("%s:%d", f.File, f.Line)).Font("monospace").FontSize(t.Rem(0.8)).SingleLine()
			})
			if b.Clicked() && !on {
				*current = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("CallStack", "selected", "index", i)
			}
		}
	})
}
