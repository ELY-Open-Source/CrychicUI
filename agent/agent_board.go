package agent

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BoardAgent is one card of an AgentBoard.
type BoardAgent struct {
	Name   string
	Status RunState
}

// AgentBoardOptions tune an AgentBoard.
type AgentBoardOptions struct {
	Status AgentStatusOptions
}

// AgentBoard shows agents as cards; cards move by status.
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
	t := c.Theme()
	rank := [4]int{2, 0, 1, 3}
	order := slices.Clone(agents)
	slices.SortStableFunc(order, func(a, b BoardAgent) int { return rank[a.Status] - rank[b.Status] })
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		for _, a := range order {
			c.Key(a.Name)
			card := kit.Card(c).Width(t.Space(40)).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(3))
			card.Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true})
			kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			card.Children(func() {
				display.Avatar(c, a.Name, nil)
				ui.Text(c, a.Name).Bold().SingleLine()
				AgentStatus(c, a.Status, o.Status)
			})
		}
	})
}
