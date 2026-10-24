package agent

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Command is one entry of SlashCommands.
type Command struct {
	Name, Description, Icon string
}

// SlashCommandsOptions word the empty list.
type SlashCommandsOptions struct {
	Empty string
}

// SlashCommands lists the commands matching a "/" prompt.
func SlashCommands(c *ui.Context, label string, commands []Command, query string, chosen *string, o SlashCommandsOptions) ui.Element {
	if label == "" || len(commands) == 0 || chosen == nil {
		panic(fmt.Sprintf("agent: SlashCommands %q needs a label, commands and a choice", label))
	}
	seen := map[string]bool{}
	for _, cmd := range commands {
		if len(cmd.Name) < 2 || cmd.Name[0] != '/' || cmd.Icon == "" || seen[cmd.Name] {
			panic(fmt.Sprintf("agent: SlashCommands %q commands need unique /names and icons", label))
		}
		seen[cmd.Name] = true
	}
	def(&o.Empty, "No matching commands")
	t := c.Theme()
	q := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(query), "/"))
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Shadow(0, 4, 12, 0, kit.Shade(0.12)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		n := 0
		for _, cmd := range commands {
			if !strings.HasPrefix(strings.ToLower(cmd.Name[1:]), q) {
				continue
			}
			c.Key(cmd.Name)
			on := cmd.Name == *chosen
			b := kit.Press(ui.ButtonBase(c).Label(cmd.Name).Row().FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).
				Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).Checked(on))
			kit.Arrive(b, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(n)*25*time.Millisecond)
			n++
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.12)
			}
			b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() {
				tile(c, cmd.Icon, t.Accent, t.Space(7))
				ui.Column(c).Grow(1).MinWidth(0).Children(func() {
					ui.Text(c, cmd.Name).Font("monospace").FontSize(t.Rem(0.85)).Bold().SingleLine()
					if cmd.Description != "" {
						kit.Muted(c, cmd.Description).FontSize(t.Rem(0.8)).MinWidth(0).SingleLine()
					}
				})
			})
			if b.Clicked() && !on {
				*chosen = cmd.Name
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SlashCommands", "selected")
			}
		}
		if n == 0 {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
				ui.Icon(c, icon.Must("search-x")).Size(16, 16).TextColor(t.TextMuted)
				kit.Muted(c, o.Empty)
			})
		}
	})
}
