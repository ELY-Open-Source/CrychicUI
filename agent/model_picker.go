package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Model is one choice of a ModelPicker.
type Model struct {
	Name, Speed, Cost string
}

// ModelPicker lists models with speed and cost tags.
func ModelPicker(c *ui.Context, label string, models []Model, selected *int) ui.Element {
	if label == "" || len(models) == 0 || selected == nil || *selected < -1 || *selected >= len(models) {
		panic(fmt.Sprintf("agent: ModelPicker %q needs a label, models and a choice or -1", label))
	}
	for _, m := range models {
		if m.Name == "" {
			panic(fmt.Sprintf("agent: ModelPicker %q models need names", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, m := range models {
			c.Key(i)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(m.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Checked(on))
			kit.Arrive(b, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
				ui.Text(c, m.Name).Bold().Grow(1).SingleLine()
				if m.Speed != "" {
					tag(c, m.Speed)
				}
				if m.Cost != "" {
					tag(c, m.Cost)
				}
			})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("ModelPicker", "selected", "index", i)
			}
		}
	})
}
