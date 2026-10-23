package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlashCommands lists the commands of a "/" prompt.
func SlashCommands(c *ui.Context, label string, commands []string, chosen *int) ui.Element {
	if label == "" || len(commands) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(commands) {
		panic(fmt.Sprintf("agent: SlashCommands %q needs a label, commands and a choice or -1", label))
	}
	for _, cmd := range commands {
		if len(cmd) < 2 || cmd[0] != '/' {
			panic(fmt.Sprintf("agent: SlashCommands %q commands start with /", label))
		}
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(56)).Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Radius).
		Background(t.Background).Border(1, t.Border).Shadow(0, 4, 12, 0, kit.Shade(0.12)).FocusGroup(ui.Vertical).Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		for i, cmd := range commands {
			c.Key(i)
			on := i == *chosen
			b := kit.Press(ui.ButtonBase(c).Label(cmd).Row().FillWidth().Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.12)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
				ui.Text(c, cmd).Font("monospace").FontSize(t.Rem(0.9)).SingleLine()
			})
			if b.Clicked() && !on {
				*chosen = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SlashCommands", "selected", "index", i)
			}
		}
	})
}
