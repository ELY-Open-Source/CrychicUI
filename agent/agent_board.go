package agent

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BoardAgent is one card of an AgentBoard.
type BoardAgent struct {
	Name, Task string
	Status     RunState
}

// AgentBoardOptions tune an AgentBoard.
type AgentBoardOptions struct {
	Status AgentStatusOptions
	Empty  string
}

// AgentBoard shows agents as cards in status columns; cards move.
func AgentBoard(c *ui.Context, label string, agents []BoardAgent, o AgentBoardOptions) ui.Element {
	if label == "" || len(agents) == 0 {
		panic(fmt.Sprintf("agent: AgentBoard %q needs a label and agents", label))
	}
	seen := map[string]bool{}
	for _, a := range agents {
		if a.Name == "" || a.Status >= runStateCount || seen[a.Name] {
			panic(fmt.Sprintf("agent: AgentBoard %q needs unique named agents in known states", label))
		}
		seen[a.Name] = true
	}
	def(&o.Status.Idle, "Idle")
	def(&o.Status.Working, "Working")
	def(&o.Status.Blocked, "Blocked")
	def(&o.Status.Done, "Done")
	def(&o.Empty, "None")
	t := c.Theme()
	names := [4]string{o.Status.Idle, o.Status.Working, o.Status.Blocked, o.Status.Done}
	root := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	board := root.ID()
	// Cards arrive once in their first column; moves fly.
	born := ui.Local(root, "born", func() time.Time { return c.Now() })
	first := ui.Local(root, "first", func() map[string]RunState {
		m := map[string]RunState{}
		for _, a := range agents {
			m[a.Name] = a.Status
		}
		return m
	})
	young := c.Now().Sub(*born) < time.Second
	if young {
		c.After(time.Second - c.Now().Sub(*born))
	}
	return root.Children(func() {
		for _, state := range []RunState{Working, Blocked, Idle, Done} {
			c.Key(state)
			n := 0
			for _, a := range agents {
				if a.Status == state {
					n++
				}
			}
			ui.Column(c).MinWidth(0).Gap(t.Space(2)).Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface).
				Role(ui.RoleList).Label(names[state]).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					AgentStatus(c, state, o.Status)
					ui.Box(c).Grow(1)
					kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
				})
				if n == 0 {
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Padding(t.Space(2))
				}
				k := 0
				for _, a := range agents {
					if a.Status != state {
						continue
					}
					c.Key(a.Name)
					card := kit.Card(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2)).Role(ui.RoleListItem).Label(a.Name).
						Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true, Match: boardCard{board, a.Name}})
					if s, ok := (*first)[a.Name]; ok && s != a.Status {
						delete(*first, a.Name)
					}
					if _, ok := (*first)[a.Name]; ok && young {
						kit.Arrive(card, ui.Pose{Y: t.Space(2), Scale: 0.96, Fade: 1}, time.Duration(k)*50*time.Millisecond)
					}
					k++
					kit.Lift(card)
					card.Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							AgentAvatar(c, a.Name, nil, a.Status == Working)
							ui.Text(c, a.Name).Bold().Grow(1).MinWidth(0).SingleLine()
						})
						if a.Task != "" {
							kit.Muted(c, a.Task).FontSize(t.Rem(0.8)).Wrap()
						}
					})
				}
			})
		}
	})
}

// boardCard flies a card between status columns.
type boardCard struct {
	board uint64
	name  string
}
