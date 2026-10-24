package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Capability is one thing an agent can do.
type Capability struct {
	Icon, Title, Text string
}

// AgentWelcomeOptions word the subtitle.
type AgentWelcomeOptions struct {
	Subtitle string
}

// AgentWelcome greets with starter prompts and capabilities; rises in.
func AgentWelcome(c *ui.Context, greeting string, prompts []Suggestion, capabilities []Capability, chosen *int, o AgentWelcomeOptions) ui.Element {
	if greeting == "" || len(prompts) == 0 || chosen == nil || *chosen < -1 || *chosen >= len(prompts) {
		panic(fmt.Sprintf("agent: AgentWelcome %q needs a greeting, prompts and a choice or -1", greeting))
	}
	seen := map[string]bool{}
	for _, p := range prompts {
		if p.Text == "" || p.Icon == "" || seen[p.Text] {
			panic(fmt.Sprintf("agent: AgentWelcome %q prompts need unique text and an icon", greeting))
		}
		seen[p.Text] = true
	}
	titles := map[string]bool{}
	for _, k := range capabilities {
		if k.Icon == "" || k.Title == "" || titles[k.Title] {
			panic(fmt.Sprintf("agent: AgentWelcome %q capabilities need unique titles and icons", greeting))
		}
		titles[k.Title] = true
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(5)).AlignItems(ui.Center).Padding(t.Space(6), t.Space(4)).Role(ui.RoleGroup).Label(greeting)
	kit.Arrive(root, ui.Pose{Y: t.Space(4), Fade: 1}, 0)
	return root.Children(func() {
		ui.Column(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			orb := ui.Box(c).Size(t.Space(14), t.Space(14)).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText)
			kit.Arrive(orb, ui.Pose{Scale: 0.6, Fade: 1}, 80*time.Millisecond).Children(func() {
				ui.Icon(c, icon.Must("sparkles")).Size(t.Space(7), t.Space(7))
			})
			ui.Text(c, greeting).Bold().FontSize(t.Rem(1.4)).Center()
			if o.Subtitle != "" {
				kit.Muted(c, o.Subtitle).Center().Wrap()
			}
		})
		if len(capabilities) > 0 {
			ui.Grid(c).Columns(len(capabilities)).FillWidth().Gap(t.Space(2)).Children(func() {
				for i, k := range capabilities {
					c.Key(k.Title)
					card := ui.Column(c).MinWidth(0).Gap(t.Space(1)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Surface)
					kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 150*time.Millisecond+time.Duration(i)*60*time.Millisecond)
					card.Children(func() {
						tile(c, k.Icon, t.Accent, t.Space(7))
						ui.Text(c, k.Title).Bold().FontSize(t.Rem(0.9)).Wrap()
						if k.Text != "" {
							kit.Muted(c, k.Text).FontSize(t.Rem(0.8)).Wrap()
						}
					})
				}
			})
		}
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Justify(ui.Center).FocusGroup(ui.Horizontal).Children(func() {
			for i, p := range prompts {
				c.Key(p.Text)
				on := i == *chosen
				b := kit.Press(ui.ButtonBase(c).Label(p.Text).MaxWidthPercent(100).
					Padding(t.Space(1.5), t.Space(3)).Radius(999).Checked(on))
				kit.Lift(b)
				kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, 300*time.Millisecond+time.Duration(i)*50*time.Millisecond)
				b.Border(1, b.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on || b.Hovered())), ui.Snappy))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.08), unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Icon(c, icon.Must(p.Icon)).Size(16, 16).TextColor(t.Accent)
						ui.Text(c, p.Text).FontSize(t.Rem(0.9)).Shrink(0).MaxWidth(t.Space(44)).Wrap()
					})
				})
				if b.Clicked() && !on {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("AgentWelcome", "picked", "index", i)
				}
			}
		})
	})
}
