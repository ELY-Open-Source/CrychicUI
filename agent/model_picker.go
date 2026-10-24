package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Model is one choice of a ModelPicker.
type Model struct {
	Name, Description, Speed, Cost string
}

// ModelPicker lists models with speed and cost tags.
func ModelPicker(c *ui.Context, label string, models []Model, selected *int) ui.Element {
	if label == "" || len(models) == 0 || selected == nil || *selected < -1 || *selected >= len(models) {
		panic(fmt.Sprintf("agent: ModelPicker %q needs a label, models and a choice or -1", label))
	}
	seen := map[string]bool{}
	for _, m := range models {
		if m.Name == "" || seen[m.Name] {
			panic(fmt.Sprintf("agent: ModelPicker %q models need unique names", label))
		}
		seen[m.Name] = true
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, m := range models {
			c.Key(m.Name)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(m.Name).Row().FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).
				Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).Checked(on))
			kit.Arrive(b, ui.Pose{Y: t.Space(1.5), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				ink := t.TextMuted
				if on {
					ink = t.Accent
				}
				tile(c, "sparkles", ink, t.Space(8))
				ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, m.Name).Bold().SingleLine()
						if m.Speed != "" {
							tag(c, m.Speed)
						}
						if m.Cost != "" {
							tag(c, m.Cost)
						}
					})
					if m.Description != "" {
						kit.Muted(c, m.Description).FontSize(t.Rem(0.8)).SingleLine()
					}
				})
				check := ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Accent)
				check.Opacity(check.AnimateTo("on", unit(on), ui.Snappy))
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
