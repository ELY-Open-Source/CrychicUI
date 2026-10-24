package agent

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PromptHistoryOptions name the list and buttons.
type PromptHistoryOptions struct {
	Label, Older, Newer, Format string
}

// PromptHistory recalls earlier prompts with Up and Down; cue Tick.
func PromptHistory(c *ui.Context, entries []string, at *int, o PromptHistoryOptions) ui.Element {
	if len(entries) == 0 || at == nil || *at < 0 || *at >= len(entries) {
		panic(fmt.Sprintf("agent: PromptHistory needs entries and a position, got %d of %d", deref(at), len(entries)))
	}
	for _, e := range entries {
		if e == "" {
			panic("agent: PromptHistory entries must not be empty")
		}
	}
	def(&o.Label, "Recent prompts")
	def(&o.Older, "Older")
	def(&o.Newer, "Newer")
	def(&o.Format, "%d of %d")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Background).
		Focusable().Role(ui.RoleList).Label(o.Label)
	root.Border(1, root.AnimateColor("ring", t.Border.Mix(t.Accent, unit(root.FocusWithin())), ui.Snappy))
	var older, newer ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(0, 0, 0, t.Space(1)).Children(func() {
			ui.Icon(c, icon.Must("history")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, o.Label).Bold().FontSize(t.Rem(0.85)).Grow(1).SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Format, *at+1, len(entries))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			older = round(c, o.Older, "arrow-up").Disabled(*at == len(entries)-1)
			newer = round(c, o.Newer, "arrow-down").Disabled(*at == 0)
		})
		for i, e := range entries {
			c.Key(i)
			on := i == *at
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius).
				Role(ui.RoleListItem).Label(e)
			kit.Arrive(row, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*30*time.Millisecond)
			row.Background(row.AnimateColor("bg", t.Accent.Alpha(0.12*unit(on)), ui.Snappy))
			row.Children(func() {
				bar := ui.Box(c).Width(t.Space(0.5)).Height(t.Space(4)).Radius(999)
				bar.Background(bar.AnimateColor("bar", t.Accent.Alpha(unit(on)), ui.Snappy))
				text := ui.Text(c, e).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
				if !on {
					text.TextColor(t.TextMuted)
				}
			})
		}
	})
	move := func(d int) {
		*at += d
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("PromptHistory", "at", "index", *at)
	}
	switch {
	case (older.Clicked() || root.Shortcut(0, ui.KeyUp)) && *at < len(entries)-1:
		move(1)
	case (newer.Clicked() || root.Shortcut(0, ui.KeyDown)) && *at > 0:
		move(-1)
	}
	return root
}

func deref(n *int) int {
	if n == nil {
		return -1
	}
	return *n
}
