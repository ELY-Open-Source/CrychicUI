package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PromptHistoryOptions name the buttons.
type PromptHistoryOptions struct {
	Older, Newer string
	Format       string
}

// PromptHistory walks earlier prompts up and down.
func PromptHistory(c *ui.Context, entries []string, at *int, o PromptHistoryOptions) ui.Element {
	if len(entries) == 0 || at == nil || *at < 0 || *at >= len(entries) {
		panic(fmt.Sprintf("agent: PromptHistory needs entries and a position, got %d", *at))
	}
	for _, e := range entries {
		if e == "" {
			panic("agent: PromptHistory entries must not be empty")
		}
	}
	def(&o.Older, "Older")
	def(&o.Newer, "Newer")
	def(&o.Format, "%d of %d")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)).Role(ui.RoleGroup).Label(fmt.Sprintf(o.Format, *at+1, len(entries)))
	var older, newer ui.Element
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			older = round(c, o.Older, "arrow-up").Disabled(*at == len(entries)-1)
			kit.Numeric(c, fmt.Sprintf(o.Format, *at+1, len(entries))).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
			newer = round(c, o.Newer, "arrow-down").Disabled(*at == 0)
		})
		kit.Muted(c, entries[*at]).SingleLine().FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
	})
	move := func(d int) {
		*at += d
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("PromptHistory", "at", "index", *at)
	}
	if older.Clicked() && *at < len(entries)-1 {
		move(1)
	}
	if newer.Clicked() && *at > 0 {
		move(-1)
	}
	return root
}
