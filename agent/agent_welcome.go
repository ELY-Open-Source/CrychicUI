package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AgentWelcome greets with starter prompts and capabilities.
func AgentWelcome(c *ui.Context, greeting string, prompts []string, capabilities []string, chosen *int) ui.Element {
	if greeting == "" || len(prompts) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(prompts) {
		panic(fmt.Sprintf("agent: AgentWelcome %q needs a greeting, prompts and a choice or -1", greeting))
	}
	for _, p := range prompts {
		if p == "" {
			panic(fmt.Sprintf("agent: AgentWelcome %q prompts must not be empty", greeting))
		}
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(8)).Role(ui.RoleGroup).Label(greeting)
	kit.Arrive(root, ui.Pose{Y: t.Space(4), Fade: 1}, 0)
	return root.Children(func() {
		ui.Icon(c, icon.Must("sparkles")).Size(t.Space(8), t.Space(8)).TextColor(t.Accent)
		ui.Text(c, greeting).Bold().FontSize(t.Rem(1.3))
		if len(capabilities) > 0 {
			ui.Row(c).Wrap().Gap(t.Space(3)).Justify(ui.Center).Children(func() {
				for _, cap := range capabilities {
					if cap == "" {
						panic(fmt.Sprintf("agent: AgentWelcome %q capabilities must not be empty", greeting))
					}
					c.Key(cap)
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success)
						kit.Muted(c, cap).FontSize(t.Rem(0.85)).SingleLine()
					})
				}
			})
		}
		ui.Row(c).Wrap().Gap(t.Space(1.5)).Justify(ui.Center).FocusGroup(ui.Horizontal).Children(func() {
			for i, p := range prompts {
				c.Key(i)
				b := kit.Press(ui.ButtonBase(c).Label(p).Padding(t.Space(1.5), t.Space(3)).Radius(999).Border(1, t.Border))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Children(func() { ui.Text(c, p).SingleLine() })
				if b.Clicked() {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("AgentWelcome", "picked", "index", i)
				}
			}
		})
	})
}
